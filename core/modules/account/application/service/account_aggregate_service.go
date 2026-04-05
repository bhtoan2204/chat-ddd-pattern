package service

import (
	"context"
	"errors"
	"time"

	"go-socket/core/modules/account/domain/aggregate"
	"go-socket/core/modules/account/domain/entity"
	"go-socket/core/modules/account/domain/repos"
	eventpkg "go-socket/core/shared/pkg/event"
	"go-socket/core/shared/pkg/stackErr"
)

type AccountAggregateService struct{}

func NewAccountAggregateService() *AccountAggregateService {
	return &AccountAggregateService{}
}

func (s *AccountAggregateService) PublishAccountCreated(ctx context.Context, outboxRepo repos.AccountOutboxEventsRepository, account *entity.Account) error {
	accountAggregate, err := aggregate.NewAccountAggregate(account.ID)
	if err != nil {
		return stackErr.Error(err)
	}
	if err := accountAggregate.RecordAccountCreated(account.Email.Value(), account.DisplayName, account.Status, account.CreatedAt); err != nil {
		return stackErr.Error(err)
	}
	return eventpkg.NewPublisher(outboxRepo).PublishAggregate(ctx, accountAggregate)
}

func (s *AccountAggregateService) PublishProfileUpdated(ctx context.Context, outboxRepo repos.AccountOutboxEventsRepository, account *entity.Account) error {
	accountAggregate, err := aggregate.NewAccountAggregate(account.ID)
	if err != nil {
		return stackErr.Error(err)
	}
	if err := accountAggregate.RecordProfileUpdated(account.DisplayName, account.Username, account.AvatarObjectKey, account.UpdatedAt); err != nil {
		return stackErr.Error(err)
	}
	return eventpkg.NewPublisher(outboxRepo).PublishAggregate(ctx, accountAggregate)
}

func (s *AccountAggregateService) PublishEmailVerificationRequested(ctx context.Context, outboxRepo repos.AccountOutboxEventsRepository, account *entity.Account, token string, requestedAt time.Time) error {
	accountAggregate, err := aggregate.NewAccountAggregate(account.ID)
	if err != nil {
		return stackErr.Error(err)
	}
	if err := accountAggregate.RecordEmailVerificationRequested(account.Email.Value(), token, requestedAt); err != nil {
		return stackErr.Error(err)
	}
	return eventpkg.NewPublisher(outboxRepo).PublishAggregate(ctx, accountAggregate)
}

func (s *AccountAggregateService) PublishEmailVerified(ctx context.Context, outboxRepo repos.AccountOutboxEventsRepository, account *entity.Account) error {
	if account.EmailVerifiedAt == nil {
		return stackErr.Error(errors.New("email_verified_at is required"))
	}

	accountAggregate, err := aggregate.NewAccountAggregate(account.ID)
	if err != nil {
		return stackErr.Error(err)
	}
	if err := accountAggregate.RecordEmailVerified(*account.EmailVerifiedAt); err != nil {
		return stackErr.Error(err)
	}
	return eventpkg.NewPublisher(outboxRepo).PublishAggregate(ctx, accountAggregate)
}

func (s *AccountAggregateService) PublishPasswordChanged(ctx context.Context, outboxRepo repos.AccountOutboxEventsRepository, account *entity.Account) error {
	if account.PasswordChangedAt == nil {
		return stackErr.Error(errors.New("password_changed_at is required"))
	}

	accountAggregate, err := aggregate.NewAccountAggregate(account.ID)
	if err != nil {
		return stackErr.Error(err)
	}
	if err := accountAggregate.RecordPasswordChanged(*account.PasswordChangedAt); err != nil {
		return stackErr.Error(err)
	}
	return eventpkg.NewPublisher(outboxRepo).PublishAggregate(ctx, accountAggregate)
}
