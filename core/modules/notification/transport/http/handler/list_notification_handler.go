package handler

import (
	"errors"
	"go-socket/core/modules/notification/application/dto/in"
	"go-socket/core/modules/notification/application/query"
	"go-socket/core/shared/pkg/logging"
	stackerr "go-socket/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type listNotificationHandler struct {
	queryBus query.Bus
}

func NewListNotificationHandler(queryBus query.Bus) *listNotificationHandler {
	return &listNotificationHandler{queryBus: queryBus}
}

func (h *listNotificationHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.ListNotificationRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, stackerr.Error(err)
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, stackerr.Error(errors.New("validate request failed"))
	}
	result, err := h.queryBus.ListNotification.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("ListNotification failed", zap.Error(err))
		return nil, stackerr.Error(errors.New("ListNotification failed"))
	}
	return result, nil
}
