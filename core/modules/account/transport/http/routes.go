package http

import (
	"go-socket/core/modules/account/application/usecase"
	"go-socket/core/modules/account/transport/http/handler"
	"go-socket/core/shared/transport/httpx"

	"github.com/gin-gonic/gin"
)

func RegisterPublicRoutes(routes *gin.RouterGroup, authUsecase usecase.AuthUsecase) {
	routes.POST("/auth/login", httpx.Wrap(handler.NewLoginHandler(authUsecase)))
	routes.POST("/auth/register", httpx.Wrap(handler.NewRegisterHandler(authUsecase)))
}

func RegisterPrivateRoutes(routes *gin.RouterGroup, authUsecase usecase.AuthUsecase) {
	routes.POST("/auth/logout", httpx.Wrap(handler.NewLogoutHandler(authUsecase)))
	routes.GET("/auth/profile", httpx.Wrap(handler.NewGetProfileHandler(authUsecase)))
}
