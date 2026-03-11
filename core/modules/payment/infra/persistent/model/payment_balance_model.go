package model

import "time"

type BalanceModel struct {
	ID        string    `gorm:"primaryKey"`
	AccountID string    `gorm:"not null;index"`
	Amount    int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (BalanceModel) TableName() string {
	return "payment_balances"
}
