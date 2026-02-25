package query

import (
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
	"go-socket/core/shared/pkg/cqrs"
)

type GetRoomHandler = cqrs.Handler[*in.GetRoomRequest, *out.GetRoomResponse]
type ListRoomHandler = cqrs.Handler[*in.ListRoomsRequest, *out.ListRoomsResponse]

type Bus struct {
	GetRoom  cqrs.Dispatcher[*in.GetRoomRequest, *out.GetRoomResponse]
	ListRoom cqrs.Dispatcher[*in.ListRoomsRequest, *out.ListRoomsResponse]
}

func NewBus(
	getRoomHandler GetRoomHandler,
	listRoomHandler ListRoomHandler,
) Bus {
	return Bus{
		GetRoom:  cqrs.NewDispatcher(getRoomHandler),
		ListRoom: cqrs.NewDispatcher(listRoomHandler),
	}
}
