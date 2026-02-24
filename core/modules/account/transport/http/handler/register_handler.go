// CODE_GENERATOR: handler
package handler

import (
	"errors"
	"go-socket/core/modules/account/application/command"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/shared/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type registerHandler struct {
	commandBus command.Bus
}

func NewRegisterHandler(commandBus command.Bus) *registerHandler {
	return &registerHandler{
		commandBus: commandBus,
	}
}

func (h *registerHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, err
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, errors.New("validate request failed")
	}
	result, err := h.commandBus.Register.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("Register failed", zap.Error(err))
		return nil, errors.New("Register failed")
	}
	return result, nil
}
