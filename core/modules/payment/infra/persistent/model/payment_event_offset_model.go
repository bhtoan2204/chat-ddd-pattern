package model

import "time"

type PaymentEventOffsetModel struct {
	ConsumerName string    `gorm:"primaryKey"`
	LastEventID  int64     `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

func (PaymentEventOffsetModel) TableName() string {
	return "payment_event_offsets"
}
