package usecase

import (
	"context"
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
)

type RoomUsecase interface {
	CreateRoom(ctx context.Context, in *in.CreateRoomRequest) (*out.CreateRoomResponse, error)
	ListRooms(ctx context.Context, in *in.ListRoomsRequest) (*out.ListRoomsResponse, error)
	GetRoom(ctx context.Context, in *in.GetRoomRequest) (*out.GetRoomResponse, error)
	UpdateRoom(ctx context.Context, in *in.UpdateRoomRequest) (*out.UpdateRoomResponse, error)
	DeleteRoom(ctx context.Context, in *in.DeleteRoomRequest) (*out.DeleteRoomResponse, error)
}
