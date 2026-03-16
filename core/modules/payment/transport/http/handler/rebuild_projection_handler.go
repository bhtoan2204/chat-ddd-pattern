package handler

import (
	"net/http"

	"go-socket/core/modules/payment/application/command"
	"go-socket/core/modules/payment/application/dto/in"
	"go-socket/core/shared/pkg/logging"
	stackerr "go-socket/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type rebuildProjectionHandler struct {
	commandBus command.Bus
}

func NewRebuildProjectionHandler(commandBus command.Bus) *rebuildProjectionHandler {
	return &rebuildProjectionHandler{commandBus: commandBus}
}

func (h *rebuildProjectionHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)

	var request in.RebuildProjectionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, nil
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, nil
	}

	result, err := h.commandBus.RebuildProjection.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("Projection rebuild failed", zap.Error(err))
		return nil, stackerr.Error(err)
	}

	return result, nil
}
