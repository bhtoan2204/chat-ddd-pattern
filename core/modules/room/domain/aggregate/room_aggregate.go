package aggregate

import (
	"go-socket/core/modules/room/domain/entity"
	"go-socket/core/modules/room/types"
)

type RoomAggregate struct {
	RoomID   string
	RoomType types.RoomType
	Messages []*entity.MessageEntity
}
