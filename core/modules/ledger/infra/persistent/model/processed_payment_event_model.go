package model

import "time"

type ProcessedPaymentEventModel struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	Provider       string    `gorm:"not null"`
	IdempotencyKey string    `gorm:"not null"`
	TransactionID  string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

func (ProcessedPaymentEventModel) TableName() string {
	return "processed_payment_events"
}
