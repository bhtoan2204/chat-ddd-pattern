package socket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *Hub) Add(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[conn] = true
}

func (h *Hub) Remove(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, conn)
}

func (h *Hub) Broadcast(message []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.clients {
		conn.WriteMessage(websocket.TextMessage, message)
	}
}
