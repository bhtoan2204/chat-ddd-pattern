package service

import (
	"context"
	"errors"
	"go-socket/core/modules/account/domain/aggregate"
	"go-socket/core/modules/account/domain/repos"
	"go-socket/core/modules/account/domain/rules"

	"gorm.io/gorm"
)

type AccountService interface {
	LoadAccountAggregate(ctx context.Context, accountID string) (*aggregate.AccountAggregate, error)
}

type accountService struct {
	baseRepo repos.Repos
}

func NewAccountService(repos repos.Repos) AccountService {
	return &accountService{
		baseRepo: repos,
	}
}

func (s *accountService) LoadAccountAggregate(ctx context.Context, accountID string) (*aggregate.AccountAggregate, error) {
	accountAggregate, err := s.baseRepo.AccountAggregateRepository().Load(ctx, accountID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, rules.ErrAccountNotFound
		}
		return nil, err
	}
	return accountAggregate, nil
}
