// CODE_GENERATOR: handler
package handler

import (
	"errors"
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/usecase"
	"go-socket/core/shared/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createMessageHandler struct {
	messageUsecase usecase.MessageUsecase
}

func NewCreateMessageHandler(messageUsecase usecase.MessageUsecase) *createMessageHandler {
	return &createMessageHandler{
		messageUsecase: messageUsecase,
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
	result, err := h.messageUsecase.CreateMessage(ctx, &request)
	if err != nil {
		logger.Errorw("CreateMessage failed", zap.Error(err))
		return nil, errors.New("CreateMessage failed")
	}
	return result, nil
}
