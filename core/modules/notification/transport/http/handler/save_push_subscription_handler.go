package handler

import (
	"errors"
	"go-socket/core/modules/notification/application/command"
	"go-socket/core/modules/notification/application/dto/in"
	"go-socket/core/shared/pkg/logging"
	stackerr "go-socket/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type savePushSubscriptionHandler struct {
	commandBus command.Bus
}

func NewSavePushSubscriptionHandler(commandBus command.Bus) *savePushSubscriptionHandler {
	return &savePushSubscriptionHandler{commandBus: commandBus}
}

func (h *savePushSubscriptionHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)

	var request in.SavePushSubscriptionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, stackerr.Error(err)
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, stackerr.Error(errors.New("validate request failed"))
	}

	result, err := h.commandBus.SavePushSubscription.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("Save push subscription failed", zap.Error(err))
		return nil, stackerr.Error(errors.New("save push subscription failed"))
	}

	return result, nil
}
