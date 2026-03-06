package entity

import "time"

type PushSubscription struct {
	ID        string
	AccountID string
	Endpoint  string
	Keys      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
