package repos

import (
	"context"
	"go-socket/core/modules/room/domain/entity"
)

type RoomMemberRepository interface {
	CreateRoomMember(ctx context.Context, roomMember *entity.RoomMemberEntity) error
}
