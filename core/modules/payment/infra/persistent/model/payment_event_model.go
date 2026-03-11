package model

import "time"

type PaymentEventModel struct {
	ID            string    `gorm:"primaryKey"`
	AggregateID   string    `gorm:"not null;uniqueIndex:idx_agg_ver"`
	AggregateType string    `gorm:"not null"`
	Version       int       `gorm:"not null;uniqueIndex:idx_agg_ver"`
	EventName     string    `gorm:"not null;index"`
	EventData     string    `gorm:"type:JSON;not null"`
	Metadata      string    `gorm:"type:JSON;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (PaymentEventModel) TableName() string {
	return "payment_events"
}
