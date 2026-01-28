// CODE_GENERATOR: routing
package http

import (
	"github.com/gin-gonic/gin"
)

type RoutingHandler interface {
	RegisterPublicHandlers(routes *gin.RouterGroup)
	RegisterPrivateHandlers(routes *gin.RouterGroup)
	Handle() gin.HandlerFunc
}

func (h *routingHandler) RegisterPublicHandlers(routes *gin.RouterGroup) {
	routes.POST("/auth/login", h.Handle())
	routes.POST("/auth/register", h.Handle())
	routes.POST("/room/create", h.Handle())
	routes.GET("/room/list", h.Handle())
	routes.GET("/room/get", h.Handle())
	routes.PUT("/room/update", h.Handle())
	routes.DELETE("/room/delete", h.Handle())
}

func (h *routingHandler) RegisterPrivateHandlers(routes *gin.RouterGroup) {
	routes.POST("/auth/logout", h.Handle())
	routes.GET("/auth/profile", h.Handle())
}
