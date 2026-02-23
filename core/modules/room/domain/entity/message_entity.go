package entity

import "time"

type MessageEntity struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	SenderID  string    `json:"sender_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
