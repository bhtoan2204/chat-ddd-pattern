package queries

import (
	"context"
	"go-socket/core/modules/room/domain/entity"
)

type RoomQueryService interface {
	GetRoom(ctx context.Context, id string) (*entity.Room, error)
}
