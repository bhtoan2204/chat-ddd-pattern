package aggregate

import (
	"errors"
	"go-socket/core/shared/pkg/event"
	"time"
)

type AccountAggregate struct {
	event.AggregateRoot

	AccountID    string
	Email        string
	Password     string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	BannedReason string
	BannedUntil  *time.Time
}

func (a *AccountAggregate) RegisterEvents(register event.RegisterEventsFunc) error {
	return register(
		&EventAccountCreated{},
		&EventAccountUpdated{},
		&EventAccountBanned{},
	)
}

func (a *AccountAggregate) Transition(e event.Event) error {
	switch data := e.EventData.(type) {
	case *EventAccountCreated:
		return a.onAccountCreated(e.AggregateID, data)
	case *EventAccountUpdated:
		return a.onAccountUpdated(data)
	case *EventAccountBanned:
		return a.onAccountBanned(data)
	default:
		return errors.New("unsupported event type")
	}
}

func (a *AccountAggregate) onAccountCreated(aggregateID string, data *EventAccountCreated) error {
	a.AccountID = aggregateID
	a.Email = data.Email
	a.CreatedAt = data.CreatedAt
	return nil
}

func (a *AccountAggregate) onAccountUpdated(data *EventAccountUpdated) error {
	a.UpdatedAt = data.UpdatedAt
	return nil
}

func (a *AccountAggregate) onAccountBanned(data *EventAccountBanned) error {
	a.BannedReason = data.BanReason
	a.BannedUntil = data.BanUntil
	return nil
}
