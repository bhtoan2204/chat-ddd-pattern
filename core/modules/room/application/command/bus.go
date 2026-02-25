package command

import (
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
	"go-socket/core/shared/pkg/cqrs"
)

type JoinRoomHandler = cqrs.Handler[*in.JoinRoomRequest, *out.JoinRoomResponse]
type CreateRoomHandler = cqrs.Handler[*in.CreateRoomRequest, *out.CreateRoomResponse]
type UpdateRoomHandler = cqrs.Handler[*in.UpdateRoomRequest, *out.UpdateRoomResponse]
type DeleteRoomHandler = cqrs.Handler[*in.DeleteRoomRequest, *out.DeleteRoomResponse]

type CreateMessageHandler = cqrs.Handler[*in.CreateMessageRequest, *out.CreateMessageResponse]

type Bus struct {
	JoinRoom      cqrs.Dispatcher[*in.JoinRoomRequest, *out.JoinRoomResponse]
	CreateRoom    cqrs.Dispatcher[*in.CreateRoomRequest, *out.CreateRoomResponse]
	UpdateRoom    cqrs.Dispatcher[*in.UpdateRoomRequest, *out.UpdateRoomResponse]
	DeleteRoom    cqrs.Dispatcher[*in.DeleteRoomRequest, *out.DeleteRoomResponse]
	CreateMessage cqrs.Dispatcher[*in.CreateMessageRequest, *out.CreateMessageResponse]
}

func NewBus(
	joinRoomHandler JoinRoomHandler,
	createRoomHandler CreateRoomHandler,
	updateRoomHandler UpdateRoomHandler,
	deleteRoomHandler DeleteRoomHandler,
	createMessageHandler CreateMessageHandler,
) Bus {
	return Bus{
		JoinRoom:      cqrs.NewDispatcher(joinRoomHandler),
		CreateRoom:    cqrs.NewDispatcher(createRoomHandler),
		UpdateRoom:    cqrs.NewDispatcher(updateRoomHandler),
		DeleteRoom:    cqrs.NewDispatcher(deleteRoomHandler),
		CreateMessage: cqrs.NewDispatcher(createMessageHandler),
	}
}
