// CODE_GENERATOR: handler
package handler

import (
	"errors"
	"wechat-clone/core/modules/room/application/dto/in"
	"wechat-clone/core/modules/room/application/dto/out"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type createMessageHandler struct {
	createMessage cqrs.Dispatcher[*in.CreateMessageRequest, *out.CreateMessageResponse]
}

func NewCreateMessageHandler(createMessage cqrs.Dispatcher[*in.CreateMessageRequest, *out.CreateMessageResponse]) *createMessageHandler {
	return &createMessageHandler{
		createMessage: createMessage,
	}
}

func (h *createMessageHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.CreateMessageRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, stackErr.Error(errors.New("validate request failed"))
	}
	result, err := h.createMessage.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("CreateMessage failed", zap.Error(err))
		return nil, stackErr.Error(errors.New("CreateMessage failed"))
	}
	return result, nil
}
