package model

import "time"

type PaymentOutboxEventModel struct {
	ID            int64  `gorm:"primaryKey"`
	AggregateID   string `gorm:"index"`
	AggregateType string `gorm:"index"`
	Version       int    `gorm:"not null"`
	EventName     string
	EventData     string
	Metadata      string
	CreatedAt     time.Time
}

func (PaymentOutboxEventModel) TableName() string {
	return "payment_outbox_events"
}
