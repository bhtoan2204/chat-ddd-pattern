package projection

import (
	"context"

	"wechat-clone/core/modules/account/domain/entity"
)

type SearchProjection interface {
	SyncAccount(ctx context.Context, account *entity.Account) error
	DeleteAccount(ctx context.Context, accountID string) error
}

type SearchRepository interface {
	SearchUsers(ctx context.Context, q string, limit, offset int) ([]*entity.Account, int64, error)
}
