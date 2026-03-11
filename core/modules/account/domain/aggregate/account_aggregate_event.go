package aggregate

import "time"

type EventAccountCreated struct {
	AccountID string
	Email     string
	CreatedAt time.Time
}

type EventAccountUpdated struct {
	AccountID string
	Email     string
	UpdatedAt time.Time
}

type EventAccountBanned struct {
	AccountID string
	BanReason string
	BanUntil  *time.Time
}
