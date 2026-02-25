// CODE_GENERATOR: handler
package handler

import (
	"errors"
	"go-socket/core/modules/room/application/command"
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/shared/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type updateRoomHandler struct {
	commandBus command.Bus
}

func NewUpdateRoomHandler(commandBus command.Bus) *updateRoomHandler {
	return &updateRoomHandler{
		commandBus: commandBus,
	}
}

func (h *updateRoomHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.UpdateRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, err
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, errors.New("validate request failed")
	}
	result, err := h.commandBus.UpdateRoom.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("UpdateRoom failed", zap.Error(err))
		return nil, errors.New("UpdateRoom failed")
	}
	return result, nil
}
