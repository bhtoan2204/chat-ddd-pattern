package socket

import (
	"sync"

	"go.uber.org/zap"
)

var _ IRoom = (*Room)(nil)

type Room struct {
	id      string
	clients map[string]IClient
	mu      sync.RWMutex
	logger  *zap.SugaredLogger
}

func NewRoom(roomID string, logger *zap.SugaredLogger) *Room {
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}
	return &Room{
		id:      roomID,
		clients: make(map[string]IClient),
		logger:  logger,
	}
}

func (r *Room) GetID() string {
	return r.id
}

func (r *Room) AddClient(client IClient) {
	r.mu.Lock()
	r.clients[client.GetID()] = client
	total := len(r.clients)
	r.mu.Unlock()

	r.logger.Debugw("client joined room", "room_id", r.id, "client_id", client.GetID(), "total_clients", total)
}

func (r *Room) RemoveClient(client IClient) {
	r.mu.Lock()
	delete(r.clients, client.GetID())
	total := len(r.clients)
	r.mu.Unlock()

	r.logger.Debugw("client left room", "room_id", r.id, "client_id", client.GetID(), "total_clients", total)
}

func (r *Room) Broadcast(message []byte) {
	r.mu.RLock()
	localClients := make([]IClient, 0, len(r.clients))
	for _, client := range r.clients {
		localClients = append(localClients, client)
	}
	r.mu.RUnlock()

	for _, client := range localClients {
		client.Send(message)
	}
}

func (r *Room) IsEmpty() bool {
	return r.ClientCount() == 0
}

func (r *Room) ClientCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}
