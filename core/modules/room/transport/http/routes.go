package http

import (
	"go-socket/core/modules/room/application/command"
	"go-socket/core/modules/room/application/query"
	"go-socket/core/modules/room/transport/http/handler"
	roomsocket "go-socket/core/modules/room/transport/websocket"
	"go-socket/core/shared/transport/httpx"

	"github.com/gin-gonic/gin"
)

func RegisterPrivateRoutes(routes *gin.RouterGroup, commandBus command.Bus, queryBus query.Bus, roomHub roomsocket.IHub) {
	routes.POST("/room/create", httpx.Wrap(handler.NewCreateRoomHandler(commandBus)))
	routes.GET("/room/list", httpx.Wrap(handler.NewListRoomsHandler(queryBus)))
	routes.GET("/room/get", httpx.Wrap(handler.NewGetRoomHandler(queryBus)))
	routes.PUT("/room/update", httpx.Wrap(handler.NewUpdateRoomHandler(commandBus)))
	routes.DELETE("/room/delete", httpx.Wrap(handler.NewDeleteRoomHandler(commandBus)))
	routes.GET("/room/ws", roomsocket.NewWSHandler(roomHub).Handle)
	// routes.POST("/message/create", httpx.Wrap(handler.NewCreateMessageHandler(messageUsecase)))
}
