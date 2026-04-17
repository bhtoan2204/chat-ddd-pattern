package repos

import (
	"context"
	"wechat-clone/core/modules/room/domain/entity"
	"wechat-clone/core/shared/utils"
)

//go:generate mockgen -package=repos -destination=room_repo_mock.go -source=room_repo.go
type RoomRepository interface {
	CreateRoom(ctx context.Context, room *entity.Room) error
	ListRooms(ctx context.Context, options utils.QueryOptions) ([]*entity.Room, error)
	GetRoomByID(ctx context.Context, id string) (*entity.Room, error)
	GetRoomByDirectKey(ctx context.Context, directKey string) (*entity.Room, error)
	UpdateRoom(ctx context.Context, room *entity.Room) error
	DeleteRoom(ctx context.Context, id string) error
}
