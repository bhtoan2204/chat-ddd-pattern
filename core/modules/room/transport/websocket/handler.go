package socket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type SocketHandler interface {
	HandleMessage(conn *websocket.Conn, message Message)
	HandleConnection(conn *websocket.Conn)
	HandleDisconnection(conn *websocket.Conn)
}

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type socketHandler struct {
	clients map[string]*websocket.Conn
	mu      sync.Mutex
}

func NewSocketHandler() SocketHandler {
	return &socketHandler{
		clients: make(map[string]*websocket.Conn),
	}
}

func (h *socketHandler) HandleMessage(conn *websocket.Conn, message Message) {
	h.mu.Lock()
	defer h.mu.Unlock()

	data, err := json.Marshal(message)
	if err != nil {
		log.Println("marshal error:", err)
		return
	}

	for id, c := range h.clients {
		if c == conn {
			continue
		}
		if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("write error to", id, err)
		}
	}
}

func (h *socketHandler) HandleConnection(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	id := conn.RemoteAddr().String()
	h.clients[id] = conn

	log.Println("socket connected:", id)
}

func (h *socketHandler) HandleDisconnection(conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	id := conn.RemoteAddr().String()

	if c, ok := h.clients[id]; ok {
		_ = c.Close()
		delete(h.clients, id)
		log.Println("socket disconnected:", id)
	}
}
