package events

import "time"

type AccountCreatedEvent struct {
	AccountID string
	Email     string
	CreatedAt time.Time
}

type AccountUpdatedEvent struct {
	AccountID string
	Email     string
	UpdatedAt time.Time
}

type AccountBannedEvent struct {
	AccountID string
	BanReason string
	BanUntil  *time.Time
}

func (e *AccountCreatedEvent) GetName() string {
	return "EventAccountCreated"
}

func (e *AccountCreatedEvent) GetData() interface{} {
	return e
}

func (e *AccountUpdatedEvent) GetName() string {
	return "EventAccountUpdated"
}

func (e *AccountUpdatedEvent) GetData() interface{} {
	return e
}

func (e *AccountBannedEvent) GetName() string {
	return "EventAccountBanned"
}

func (e *AccountBannedEvent) GetData() interface{} {
	return e
}
