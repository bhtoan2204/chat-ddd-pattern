package http

import (
	"go-socket/core/modules/room/application/usecase"
	"go-socket/core/modules/room/transport/http/handler"
	"go-socket/core/shared/transport/httpx"

	"github.com/gin-gonic/gin"
)

func RegisterPrivateRoutes(routes *gin.RouterGroup, roomUsecase usecase.RoomUsecase, messageUsecase usecase.MessageUsecase) {
	routes.POST("/room/create", httpx.Wrap(handler.NewCreateRoomHandler(roomUsecase)))
	routes.GET("/room/list", httpx.Wrap(handler.NewListRoomsHandler(roomUsecase)))
	routes.GET("/room/get", httpx.Wrap(handler.NewGetRoomHandler(roomUsecase)))
	routes.PUT("/room/update", httpx.Wrap(handler.NewUpdateRoomHandler(roomUsecase)))
	routes.DELETE("/room/delete", httpx.Wrap(handler.NewDeleteRoomHandler(roomUsecase)))
	routes.POST("/message/create", httpx.Wrap(handler.NewCreateMessageHandler(messageUsecase)))
}
