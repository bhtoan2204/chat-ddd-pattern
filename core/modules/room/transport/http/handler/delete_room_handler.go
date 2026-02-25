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

type deleteRoomHandler struct {
	commandBus command.Bus
}

func NewDeleteRoomHandler(commandBus command.Bus) *deleteRoomHandler {
	return &deleteRoomHandler{
		commandBus: commandBus,
	}
}

func (h *deleteRoomHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.DeleteRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, err
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, errors.New("validate request failed")
	}
	result, err := h.commandBus.DeleteRoom.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("DeleteRoom failed", zap.Error(err))
		return nil, errors.New("DeleteRoom failed")
	}
	return result, nil
}
