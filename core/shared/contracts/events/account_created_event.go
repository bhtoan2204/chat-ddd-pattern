package events

import "time"

type AccountCreatedEvent struct {
	AccountID string
	Email     string
	CreatedAt time.Time
}

func (e *AccountCreatedEvent) GetName() string {
	return "account.created"
}

func (e *AccountCreatedEvent) GetData() interface{} {
	return e
}
