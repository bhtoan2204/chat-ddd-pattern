package model

import "time"

type PaymentIntentModel struct {
	TransactionID   string `gorm:"primaryKey"`
	Provider        string `gorm:"not null"`
	ExternalRef     *string
	Amount          int64     `gorm:"not null"`
	Currency        string    `gorm:"not null"`
	DebitAccountID  string    `gorm:"not null"`
	CreditAccountID string    `gorm:"not null"`
	Status          string    `gorm:"not null"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (PaymentIntentModel) TableName() string {
	return "payment_intents"
}
