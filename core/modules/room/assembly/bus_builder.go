package assembly

import (
	appCtx "go-socket/core/context"
	"go-socket/core/modules/room/application/command"
	"go-socket/core/modules/room/application/query"
	roomrepo "go-socket/core/modules/room/infra/persistent/repository"
)

type Buses struct {
	Command command.Bus
	Query   query.Bus
}

func BuildBuses(appCtx *appCtx.AppContext) Buses {
	roomRepos := roomrepo.NewRepoImpl(appCtx)
	joinRoomHandler := command.NewJoinRoomHandler(roomRepos.RoomRepository())
	createRoomHandler := command.NewCreateRoomHandler(roomRepos)
	updateRoomHandler := command.NewUpdateRoomHandler(roomRepos.RoomRepository())
	deleteRoomHandler := command.NewDeleteRoomHandler(roomRepos.RoomRepository())
	createMessageHandler := command.NewCreateMessageHandler(roomRepos.MessageRepository())

	getRoomHandler := query.NewGetRoomHandler(roomRepos.RoomRepository())
	listRoomHandler := query.NewListRoomHandler(roomRepos.RoomRepository())
	return Buses{
		Command: command.NewBus(joinRoomHandler, createRoomHandler, updateRoomHandler, deleteRoomHandler, createMessageHandler),
		Query:   query.NewBus(getRoomHandler, listRoomHandler),
	}
}
