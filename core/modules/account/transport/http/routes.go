package http

import (
	"go-socket/core/modules/account/application/command"
	"go-socket/core/modules/account/application/query"
	"go-socket/core/modules/account/transport/http/handler"
	"go-socket/core/shared/transport/httpx"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(routes *gin.RouterGroup, commandBus command.Bus) {
	routes.POST("/auth/login", httpx.Wrap(handler.NewLoginHandler(commandBus)))
	routes.POST("/auth/register", httpx.Wrap(handler.NewRegisterHandler(commandBus)))
}

func RegisterPrivateRoutes(routes *gin.RouterGroup, commandBus command.Bus, queryBus query.Bus) {
	routes.POST("/auth/logout", httpx.Wrap(handler.NewLogoutHandler(commandBus)))
	routes.GET("/auth/profile", httpx.Wrap(handler.NewGetProfileHandler(queryBus)))
}
