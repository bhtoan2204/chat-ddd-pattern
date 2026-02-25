package socket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 16 * 1024
	sendBufferSize = 256
)

var _ IClient = (*Client)(nil)

type Client struct {
	id     string
	userID string
	conn   *websocket.Conn
	send   chan []byte
	logger *zap.SugaredLogger

	sendMu   sync.RWMutex
	closeMu  sync.Once
	isClosed bool
}

func NewClient(conn *websocket.Conn, clientID, userID string, logger *zap.SugaredLogger) *Client {
	if logger == nil {
		logger = zap.NewNop().Sugar()
	}
	if clientID == "" {
		clientID = uuid.NewString()
	}
	return &Client{
		id:     clientID,
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, sendBufferSize),
		logger: logger,
	}
}

func (c *Client) GetID() string {
	return c.id
}

func (c *Client) GetUserID() string {
	return c.userID
}

func (c *Client) Send(message []byte) {
	if len(message) == 0 {
		return
	}

	shouldClose := false
	c.sendMu.RLock()
	if !c.isClosed {
		select {
		case c.send <- message:
		default:
			shouldClose = true
		}
	}
	c.sendMu.RUnlock()

	if shouldClose {
		c.logger.Warnw("client send buffer is full, closing connection", "client_id", c.id, "user_id", c.userID)
		c.Close()
	}
}

func (c *Client) ReadPump(ctx context.Context, hub IHub) {
	defer hub.Unregister(c)

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	for {
		select {
		case <-ctx.Done():
			c.logger.Debugw("stopping read pump due to context cancellation", "client_id", c.id)
			return
		default:
		}

		_, payload, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Warnw("unexpected websocket read error", "client_id", c.id, "error", err)
			} else {
				c.logger.Infow("websocket connection closed while reading", "client_id", c.id, "error", err)
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(payload, &msg); err != nil {
			c.logger.Warnw("invalid websocket payload", "client_id", c.id, "error", err)
			continue
		}
		if msg.SenderID == "" {
			msg.SenderID = c.userID
		}

		if err := hub.HandleMessage(ctx, c, msg); err != nil {
			c.logger.Errorw("failed to handle websocket message", "client_id", c.id, "action", msg.Action, "room_id", msg.RoomID, "error", err)
		}
	}
}

func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	defer c.Close()

	for {
		select {
		case <-ctx.Done():
			c.logger.Debugw("stopping write pump due to context cancellation", "client_id", c.id)
			return

		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				c.logger.Errorw("failed to write websocket message", "client_id", c.id, "error", err)
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				c.logger.Infow("failed to write websocket ping", "client_id", c.id, "error", err)
				return
			}
		}
	}
}

func (c *Client) Close() {
	c.closeMu.Do(func() {
		c.sendMu.Lock()
		c.isClosed = true
		close(c.send)
		c.sendMu.Unlock()

		if err := c.conn.Close(); err != nil {
			c.logger.Debugw("error while closing websocket connection", "client_id", c.id, "error", err)
		}
	})
}
