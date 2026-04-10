package repos

import (
	"context"
	"go-socket/core/modules/room/domain/entity"
)

type RoomAccountProjectionRepository interface {
	ProjectAccount(context.Context, *entity.AccountEntity) error
	ListByAccountIDs(ctx context.Context, accountIDs []string) ([]*entity.AccountEntity, error)
}
