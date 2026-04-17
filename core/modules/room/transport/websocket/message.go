package socket

import "encoding/json"

const (
	ActionJoinRoom      = "JOIN_ROOM"
	ActionJoinRoomOK    = "JOIN_ROOM_OK"
	ActionJoinRoomError = "JOIN_ROOM_ERROR"
	ActionLeaveRoom     = "LEAVE_ROOM"
	ActionChatMessage   = "CHAT_MESSAGE"
	ActionTyping        = "TYPING"
	ActionPresence      = "PRESENCE"
	ActionSeen          = "SEEN"
)

type Message struct {
	Action       string          `json:"action"`
	RoomID       string          `json:"room_id,omitempty"`
	SenderID     string          `json:"sender_id,omitempty"`
	Data         json.RawMessage `json:"data,omitempty"`
	RecipientIDs []string        `json:"recipient_ids,omitempty"`
}

type AckMessage struct {
	Message
	IsSuccess bool   `json:"is_success"`
	Error     string `json:"error,omitempty"`
}
