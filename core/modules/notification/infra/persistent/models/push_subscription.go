package models

import "time"

type PushSubscriptionModel struct {
	ID        string    `gorm:"primaryKey;type:char(36)"`
	AccountID string    `gorm:"not null;index;index:idx_push_subscriptions_account_endpoint,unique"`
	Endpoint  string    `gorm:"not null;index:idx_push_subscriptions_account_endpoint,unique"`
	Keys      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (PushSubscriptionModel) TableName() string {
	return "push_subscriptions"
}
