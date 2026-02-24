package models

import "time"

type AccountOutboxEventModel struct {
	ID            int64     `gorm:"primaryKey;autoIncrement"`
	AggregateID   string    `gorm:"not null;index"`
	AggregateType string    `gorm:"not null;index"`
	Version       int       `gorm:"not null"`
	EventName     string    `gorm:"not null;index"`
	EventData     string    `gorm:"type:CLOB;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (AccountOutboxEventModel) TableName() string {
	return "account_outbox_events"
}
