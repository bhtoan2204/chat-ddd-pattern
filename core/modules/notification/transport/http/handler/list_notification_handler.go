package handler

import (
	"errors"
	"wechat-clone/core/modules/notification/application/dto/in"
	"wechat-clone/core/modules/notification/application/dto/out"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type listNotificationHandler struct {
	listNotification cqrs.Dispatcher[*in.ListNotificationRequest, *out.ListNotificationResponse]
}

func NewListNotificationHandler(listNotification cqrs.Dispatcher[*in.ListNotificationRequest, *out.ListNotificationResponse]) *listNotificationHandler {
	return &listNotificationHandler{listNotification: listNotification}
}

func (h *listNotificationHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.ListNotificationRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, stackErr.Error(errors.New("validate request failed"))
	}
	result, err := h.listNotification.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("ListNotification failed", zap.Error(err))
		return nil, stackErr.Error(errors.New("ListNotification failed"))
	}
	return result, nil
}
