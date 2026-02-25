package socket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const roomChannelPrefix = "room:"

var _ IHub = (*Hub)(nil)

type roomSubscription struct {
	pubsub *redis.PubSub
	cancel context.CancelFunc
	once   sync.Once
}

func (s *roomSubscription) Close() error {
	var closeErr error
	s.once.Do(func() {
		s.cancel()
		closeErr = s.pubsub.Close()
	})
	return closeErr
}

type Hub struct {
	redisClient *redis.Client
	logger      *zap.SugaredLogger

	mu            sync.RWMutex
	clients       map[string]IClient
	rooms         map[string]IRoom
	clientRooms   map[string]map[string]struct{}
	subscriptions map[string]*roomSubscription

	ctx      context.Context
	cancel   context.CancelFunc
	closeMu  sync.Once
	isClosed bool
}

func NewHub(redisClient *redis.Client, logger *zap.SugaredLogger) *Hub {
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		redisClient:   redisClient,
		logger:        logger,
		clients:       make(map[string]IClient),
		rooms:         make(map[string]IRoom),
		clientRooms:   make(map[string]map[string]struct{}),
		subscriptions: make(map[string]*roomSubscription),
		ctx:           ctx,
		cancel:        cancel,
	}
}

func (h *Hub) Register(client IClient) {
	if client == nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.isClosed {
		client.Close()
		return
	}
	h.clients[client.GetID()] = client
	if _, ok := h.clientRooms[client.GetID()]; !ok {
		h.clientRooms[client.GetID()] = make(map[string]struct{})
	}

	h.logger.Infow("client registered", "client_id", client.GetID(), "user_id", client.GetUserID(), "clients", len(h.clients))
}

func (h *Hub) Unregister(client IClient) {
	if client == nil {
		return
	}

	clientID := client.GetID()

	h.mu.Lock()
	if _, exists := h.clients[clientID]; !exists {
		h.mu.Unlock()
		client.Close()
		return
	}
	delete(h.clients, clientID)
	roomIDs := make([]string, 0, len(h.clientRooms[clientID]))
	for roomID := range h.clientRooms[clientID] {
		roomIDs = append(roomIDs, roomID)
	}
	h.mu.Unlock()

	for _, roomID := range roomIDs {
		if err := h.LeaveRoom(context.Background(), client, roomID); err != nil {
			h.logger.Warnw("failed to leave room while unregistering client", "client_id", clientID, "room_id", roomID, "error", err)
		}
	}

	h.mu.Lock()
	delete(h.clientRooms, clientID)
	remainingClients := len(h.clients)
	h.mu.Unlock()

	client.Close()
	h.logger.Infow("client unregistered", "client_id", clientID, "clients", remainingClients)
}

func (h *Hub) JoinRoom(ctx context.Context, client IClient, roomID string) error {
	if client == nil {
		return errors.New("client is nil")
	}
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return errors.New("room_id is required")
	}

	h.mu.Lock()
	if h.isClosed {
		h.mu.Unlock()
		return errors.New("hub is closed")
	}
	if _, ok := h.clients[client.GetID()]; !ok {
		h.clients[client.GetID()] = client
	}
	room, ok := h.rooms[roomID]
	if !ok {
		room = NewRoom(roomID, h.logger)
		h.rooms[roomID] = room
	}
	room.AddClient(client)

	if _, ok := h.clientRooms[client.GetID()]; !ok {
		h.clientRooms[client.GetID()] = make(map[string]struct{})
	}
	h.clientRooms[client.GetID()][roomID] = struct{}{}

	_, hasSubscription := h.subscriptions[roomID]
	h.mu.Unlock()

	if !hasSubscription {
		if err := h.subscribeRoom(roomID); err != nil {
			return err
		}
	}

	h.logger.Infow("client joined room", "client_id", client.GetID(), "room_id", roomID)
	return nil
}

func (h *Hub) LeaveRoom(ctx context.Context, client IClient, roomID string) error {
	_ = ctx
	if client == nil {
		return errors.New("client is nil")
	}
	roomID = strings.TrimSpace(roomID)
	if roomID == "" {
		return errors.New("room_id is required")
	}

	shouldUnsubscribe := false

	h.mu.Lock()
	room, exists := h.rooms[roomID]
	if exists {
		room.RemoveClient(client)
		if room.IsEmpty() {
			delete(h.rooms, roomID)
			shouldUnsubscribe = true
		}
	}

	if rooms, ok := h.clientRooms[client.GetID()]; ok {
		delete(rooms, roomID)
		if len(rooms) == 0 {
			delete(h.clientRooms, client.GetID())
		}
	}
	h.mu.Unlock()

	if shouldUnsubscribe {
		h.unsubscribeRoom(roomID)
	}

	h.logger.Infow("client left room", "client_id", client.GetID(), "room_id", roomID)
	return nil
}

func (h *Hub) HandleMessage(ctx context.Context, client IClient, msg Message) error {
	if client == nil {
		return errors.New("client is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}

	switch msg.Action {
	case ActionJoinRoom:
		return h.JoinRoom(ctx, client, msg.RoomID)

	case ActionLeaveRoom:
		return h.LeaveRoom(ctx, client, msg.RoomID)

	case ActionChatMessage:
		roomID := strings.TrimSpace(msg.RoomID)
		if roomID == "" {
			return errors.New("room_id is required for chat message")
		}
		if h.redisClient == nil {
			return errors.New("redis client is nil")
		}
		if msg.SenderID == "" {
			msg.SenderID = client.GetUserID()
		}
		if msg.SenderID == "" {
			msg.SenderID = client.GetID()
		}

		payload, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("marshal websocket message: %w", err)
		}
		if err := h.redisClient.Publish(ctx, roomChannelName(roomID), payload).Err(); err != nil {
			return fmt.Errorf("publish redis message: %w", err)
		}
		return nil
	}

	return fmt.Errorf("unsupported websocket action: %s", msg.Action)
}

func (h *Hub) Close() {
	h.closeMu.Do(func() {
		h.cancel()

		h.mu.Lock()
		h.isClosed = true
		clients := make([]IClient, 0, len(h.clients))
		for _, client := range h.clients {
			clients = append(clients, client)
		}
		subscriptions := make([]*roomSubscription, 0, len(h.subscriptions))
		for _, sub := range h.subscriptions {
			subscriptions = append(subscriptions, sub)
		}

		h.clients = make(map[string]IClient)
		h.rooms = make(map[string]IRoom)
		h.clientRooms = make(map[string]map[string]struct{})
		h.subscriptions = make(map[string]*roomSubscription)
		h.mu.Unlock()

		for _, sub := range subscriptions {
			if err := sub.Close(); err != nil {
				h.logger.Warnw("failed to close redis pubsub", "error", err)
			}
		}
		for _, client := range clients {
			client.Close()
		}

		h.logger.Infow("hub closed")
	})
}

func (h *Hub) subscribeRoom(roomID string) error {
	if h.redisClient == nil {
		return errors.New("redis client is nil")
	}

	h.mu.Lock()
	if h.isClosed {
		h.mu.Unlock()
		return errors.New("hub is closed")
	}
	if _, exists := h.subscriptions[roomID]; exists {
		h.mu.Unlock()
		return nil
	}
	subCtx, cancel := context.WithCancel(h.ctx)
	pubsub := h.redisClient.Subscribe(subCtx, roomChannelName(roomID))
	sub := &roomSubscription{
		pubsub: pubsub,
		cancel: cancel,
	}
	h.subscriptions[roomID] = sub
	h.mu.Unlock()

	if _, err := pubsub.Receive(subCtx); err != nil {
		h.removeSubscription(roomID, sub)
		_ = sub.Close()
		return fmt.Errorf("subscribe to redis room channel: %w", err)
	}

	go h.consumeRoomMessages(roomID, sub)
	h.logger.Infow("subscribed redis room channel", "room_id", roomID, "channel", roomChannelName(roomID))
	return nil
}

func (h *Hub) unsubscribeRoom(roomID string) {
	sub := h.detachSubscription(roomID)
	if sub == nil {
		return
	}
	if err := sub.Close(); err != nil {
		h.logger.Warnw("failed to unsubscribe redis room channel", "room_id", roomID, "error", err)
		return
	}
	h.logger.Infow("unsubscribed redis room channel", "room_id", roomID, "channel", roomChannelName(roomID))
}

func (h *Hub) consumeRoomMessages(roomID string, sub *roomSubscription) {
	defer func() {
		h.removeSubscription(roomID, sub)
		if err := sub.Close(); err != nil {
			h.logger.Debugw("error while closing redis pubsub from consumer", "room_id", roomID, "error", err)
		}
	}()

	channel := sub.pubsub.Channel()
	for {
		select {
		case <-h.ctx.Done():
			return
		case message, ok := <-channel:
			if !ok {
				return
			}
			h.broadcastLocal(roomID, []byte(message.Payload))
		}
	}
}

func (h *Hub) broadcastLocal(roomID string, payload []byte) {
	h.mu.RLock()
	room, ok := h.rooms[roomID]
	h.mu.RUnlock()
	if !ok {
		return
	}
	room.Broadcast(payload)
}

func (h *Hub) detachSubscription(roomID string) *roomSubscription {
	h.mu.Lock()
	defer h.mu.Unlock()

	sub := h.subscriptions[roomID]
	delete(h.subscriptions, roomID)
	return sub
}

func (h *Hub) removeSubscription(roomID string, expected *roomSubscription) {
	h.mu.Lock()
	defer h.mu.Unlock()

	current, ok := h.subscriptions[roomID]
	if !ok {
		return
	}
	if expected != nil && current != expected {
		return
	}
	delete(h.subscriptions, roomID)
}

func roomChannelName(roomID string) string {
	return roomChannelPrefix + roomID
}
