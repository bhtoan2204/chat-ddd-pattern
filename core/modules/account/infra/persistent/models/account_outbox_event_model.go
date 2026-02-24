package models

import "time"

type AccountOutboxEventModel struct {
	ID        string      `gorm:"primaryKey,autoIncrement"`
	EventName string      `gorm:"not null"`
	EventData interface{} `gorm:"not null"`
	CreatedAt time.Time   `gorm:"autoCreateTime"`
}

func (AccountOutboxEventModel) TableName() string {
	return "account_outbox_events"
}
