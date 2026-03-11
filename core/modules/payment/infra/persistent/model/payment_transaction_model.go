package model

import "time"

type PaymentTransactionModel struct {
	ID        string    `gorm:"primaryKey"`
	AccountID string    `gorm:"not null;index"`
	EventID   string    `gorm:"not null;index"`
	Amount    int64     `gorm:"not null"`
	Type      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (PaymentTransactionModel) TableName() string {
	return "payment_transactions"
}
