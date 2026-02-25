package command

import (
	"context"
	"fmt"
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
	"go-socket/core/modules/room/domain/repos"
)

type joinRoomHandler struct {
	roomRepo repos.RoomRepository
}

func NewJoinRoomHandler(roomRepo repos.RoomRepository) JoinRoomHandler {
	return &joinRoomHandler{
		roomRepo: roomRepo,
	}
}

func (h *joinRoomHandler) Handle(ctx context.Context, req *in.JoinRoomRequest) (*out.JoinRoomResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
