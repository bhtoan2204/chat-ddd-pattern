package aggregate

import (
	"reflect"
	"time"

	"go-socket/core/shared/pkg/stackErr"
)

func NewAccountAggregate(accountID string) (*AccountAggregate, error) {
	agg := &AccountAggregate{}
	agg.SetAggregateType(reflect.TypeOf(agg).Elem().Name())
	if err := agg.SetID(accountID); err != nil {
		return nil, stackErr.Error(err)
	}

	return agg, nil
}

func (a *AccountAggregate) RecordAccountCreated(email, displayName, status string, createdAt time.Time) error {
	return a.ApplyChange(a, &EventAccountCreated{
		AccountID:   a.AggregateID(),
		Email:       email,
		DisplayName: displayName,
		Status:      status,
		CreatedAt:   createdAt,
	})
}

func (a *AccountAggregate) RecordProfileUpdated(displayName string, username, avatarObjectKey *string, updatedAt time.Time) error {
	return a.ApplyChange(a, &EventAccountProfileUpdated{
		AccountID:       a.AggregateID(),
		DisplayName:     displayName,
		Username:        username,
		AvatarObjectKey: avatarObjectKey,
		UpdatedAt:       updatedAt,
	})
}

func (a *AccountAggregate) RecordEmailVerificationRequested(email, token string, requestedAt time.Time) error {
	return a.ApplyChange(a, &EventAccountEmailVerificationRequested{
		AccountID:         a.AggregateID(),
		Email:             email,
		VerificationToken: token,
		RequestedAt:       requestedAt,
	})
}

func (a *AccountAggregate) RecordEmailVerified(verifiedAt time.Time) error {
	return a.ApplyChange(a, &EventAccountEmailVerified{
		AccountID:       a.AggregateID(),
		EmailVerifiedAt: verifiedAt,
	})
}

func (a *AccountAggregate) RecordPasswordChanged(changedAt time.Time) error {
	return a.ApplyChange(a, &EventAccountPasswordChanged{
		AccountID:         a.AggregateID(),
		PasswordChangedAt: changedAt,
	})
}
