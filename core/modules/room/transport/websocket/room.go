package socket

import (
	"context"
	"go-socket/core/shared/pkg/logging"
	"sync"
)

var _ IRoom = (*Room)(nil)

type Room struct {
	id      string
	clients map[string]IClient
	mu      sync.RWMutex
}

func NewRoom(ctx context.Context, roomID string) *Room {
	return &Room{
		id:      roomID,
		clients: make(map[string]IClient),
	}
}

func (r *Room) GetID() string {
	return r.id
}

func (r *Room) AddClient(ctx context.Context, client IClient) {
	log := logging.FromContext(ctx)
	r.mu.Lock()
	r.clients[client.GetID()] = client
	total := len(r.clients)
	r.mu.Unlock()

	log.Debugw("client joined room", "room_id", r.id, "client_id", client.GetID(), "total_clients", total)
}

func (r *Room) RemoveClient(ctx context.Context, client IClient) {
	log := logging.FromContext(ctx)
	r.mu.Lock()
	delete(r.clients, client.GetID())
	total := len(r.clients)
	r.mu.Unlock()

	log.Debugw("client left room", "room_id", r.id, "client_id", client.GetID(), "total_clients", total)
}

func (r *Room) Broadcast(ctx context.Context, message []byte) {
	r.mu.RLock()
	localClients := make([]IClient, 0, len(r.clients))
	for _, client := range r.clients {
		localClients = append(localClients, client)
	}
	r.mu.RUnlock()

	for _, client := range localClients {
		client.Send(ctx, message)
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
