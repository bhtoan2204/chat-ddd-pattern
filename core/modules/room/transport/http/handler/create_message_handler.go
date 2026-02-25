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

type createMessageHandler struct {
	commandBus command.Bus
}

func NewCreateMessageHandler(commandBus command.Bus) *createMessageHandler {
	return &createMessageHandler{
		commandBus: commandBus,
	}
}

func (h *createMessageHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.CreateMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, err
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, errors.New("validate request failed")
	}
	result, err := h.commandBus.CreateMessage.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("CreateMessage failed", zap.Error(err))
		return nil, errors.New("CreateMessage failed")
	}
	return result, nil
}
