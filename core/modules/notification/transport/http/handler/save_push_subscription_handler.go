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

type savePushSubscriptionHandler struct {
	savePushSubscription cqrs.Dispatcher[*in.SavePushSubscriptionRequest, *out.SavePushSubscriptionResponse]
}

func NewSavePushSubscriptionHandler(savePushSubscription cqrs.Dispatcher[*in.SavePushSubscriptionRequest, *out.SavePushSubscriptionResponse]) *savePushSubscriptionHandler {
	return &savePushSubscriptionHandler{savePushSubscription: savePushSubscription}
}

func (h *savePushSubscriptionHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)

	var request in.SavePushSubscriptionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, stackErr.Error(errors.New("validate request failed"))
	}

	result, err := h.savePushSubscription.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("Save push subscription failed", zap.Error(err))
		return nil, stackErr.Error(errors.New("save push subscription failed"))
	}

	return result, nil
}
