package repos

import (
	"context"
	"go-socket/core/domain/entity"
)

type AccountRepository interface {
	GetAccountByID(ctx context.Context, id string) (*entity.Account, error)
	CreateAccount(ctx context.Context, account *entity.Account) error
	UpdateAccount(ctx context.Context, account *entity.Account) error
	DeleteAccount(ctx context.Context, id string) error
}
