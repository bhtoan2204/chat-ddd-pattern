package http

import (
	"go-socket/config"
	appCtx "go-socket/core/context"
	"go-socket/core/delivery/http/handler"
	"go-socket/core/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type routingConfig struct {
	handler    handler.RequestHandler
	permission []string
}

type routingHandler struct {
	config      *config.Config
	appCtx      *appCtx.AppContext
	redisClient *redis.Client
	registry    map[string]routingConfig
}

func NewRoutingHandler(config *config.Config, appCtx *appCtx.AppContext, redisClient *redis.Client, usecase usecase.Usecase) RoutingHandler {
	return &routingHandler{
		config:      config,
		appCtx:      appCtx,
		redisClient: redisClient,
		registry:    BuildRegistry(config, usecase),
	}
}

func (h *routingHandler) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		route := c.Request.Method + ":" + c.FullPath()
		routingCfg, found := h.registry[route]
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		data, err := routingCfg.handler.Handle(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, data)
	}
}
