package models

import "time"

type MessageModel struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	RoomID    string    `gorm:"not null;index" json:"room_id"`
	SenderID  string    `gorm:"not null;index" json:"sender_id"`
	Message   string    `gorm:"type:VARCHAR2(4000);not null" json:"message"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (MessageModel) TableName() string {
	return "messages"
}
