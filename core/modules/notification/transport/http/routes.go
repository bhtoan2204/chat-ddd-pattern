package http

import (
	"go-socket/core/modules/notification/application/command"
	"go-socket/core/modules/notification/application/query"
	"go-socket/core/modules/notification/transport/http/handler"
	"go-socket/core/shared/transport/httpx"

	"github.com/gin-gonic/gin"
)

func RegisterPrivateRoutes(routes *gin.RouterGroup, commandBus command.Bus, queryBus query.Bus) {
	routes.POST("/notification/push-subscriptions", httpx.Wrap(handler.NewSavePushSubscriptionHandler(commandBus)))
	routes.GET("/notification/list", httpx.Wrap(handler.NewListNotificationHandler(queryBus)))
}
