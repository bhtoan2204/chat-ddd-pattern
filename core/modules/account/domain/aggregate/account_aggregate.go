package aggregate

import (
	"errors"
	"go-socket/core/shared/pkg/event"
	"time"
)

type AccountAggregate struct {
	event.AggregateRoot

	AccountID                      string
	Email                          string
	DisplayName                    string
	Username                       *string
	AvatarObjectKey                *string
	Status                         string
	Password                       string
	EmailVerifiedAt                *time.Time
	LastEmailVerificationRequested *time.Time
	PasswordChangedAt              *time.Time
	CreatedAt                      time.Time
	UpdatedAt                      time.Time
	BannedReason                   string
	BannedUntil                    *time.Time
}

func (a *AccountAggregate) RegisterEvents(register event.RegisterEventsFunc) error {
	return register(
		&EventAccountCreated{},
		&EventAccountUpdated{},
		&EventAccountProfileUpdated{},
		&EventAccountEmailVerificationRequested{},
		&EventAccountEmailVerified{},
		&EventAccountPasswordChanged{},
		&EventAccountBanned{},
	)
}

func (a *AccountAggregate) Transition(e event.Event) error {
	switch data := e.EventData.(type) {
	case *EventAccountCreated:
		return a.onAccountCreated(e.AggregateID, data)
	case *EventAccountUpdated:
		return a.onAccountUpdated(data)
	case *EventAccountProfileUpdated:
		return a.onAccountProfileUpdated(data)
	case *EventAccountEmailVerificationRequested:
		return a.onAccountEmailVerificationRequested(data)
	case *EventAccountEmailVerified:
		return a.onAccountEmailVerified(data)
	case *EventAccountPasswordChanged:
		return a.onAccountPasswordChanged(data)
	case *EventAccountBanned:
		return a.onAccountBanned(data)
	default:
		return errors.New("unsupported event type")
	}
}

func (a *AccountAggregate) onAccountCreated(aggregateID string, data *EventAccountCreated) error {
	a.AccountID = aggregateID
	a.Email = data.Email
	a.DisplayName = data.DisplayName
	a.Status = data.Status
	a.CreatedAt = data.CreatedAt
	a.UpdatedAt = data.CreatedAt
	return nil
}

func (a *AccountAggregate) onAccountUpdated(data *EventAccountUpdated) error {
	a.Email = data.Email
	a.UpdatedAt = data.UpdatedAt
	return nil
}

func (a *AccountAggregate) onAccountProfileUpdated(data *EventAccountProfileUpdated) error {
	a.DisplayName = data.DisplayName
	a.Username = data.Username
	a.AvatarObjectKey = data.AvatarObjectKey
	a.UpdatedAt = data.UpdatedAt
	return nil
}

func (a *AccountAggregate) onAccountEmailVerificationRequested(data *EventAccountEmailVerificationRequested) error {
	requestedAt := data.RequestedAt
	a.LastEmailVerificationRequested = &requestedAt
	a.UpdatedAt = requestedAt
	return nil
}

func (a *AccountAggregate) onAccountEmailVerified(data *EventAccountEmailVerified) error {
	verifiedAt := data.EmailVerifiedAt
	a.EmailVerifiedAt = &verifiedAt
	a.UpdatedAt = verifiedAt
	return nil
}

func (a *AccountAggregate) onAccountPasswordChanged(data *EventAccountPasswordChanged) error {
	changedAt := data.PasswordChangedAt
	a.PasswordChangedAt = &changedAt
	a.UpdatedAt = changedAt
	return nil
}

func (a *AccountAggregate) onAccountBanned(data *EventAccountBanned) error {
	a.BannedReason = data.BanReason
	a.BannedUntil = data.BanUntil
	return nil
}
