// CODE_GENERATOR - do not edit: handler
package handler

import (
	"net/http"

	"wechat-clone/core/modules/relationship/application/dto/in"
	"wechat-clone/core/modules/relationship/application/dto/out"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type rejectFriendRequestHandler struct {
	rejectFriendRequest cqrs.Dispatcher[*in.RejectFriendRequestRequest, *out.RejectFriendRequestResponse]
}

func NewRejectFriendRequestHandler(
	rejectFriendRequest cqrs.Dispatcher[*in.RejectFriendRequestRequest, *out.RejectFriendRequestResponse],
) *rejectFriendRequestHandler {
	return &rejectFriendRequestHandler{
		rejectFriendRequest: rejectFriendRequest,
	}
}

func (h *rejectFriendRequestHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.RejectFriendRequestRequest
	request.RequesterUserID = c.Param("requester_user_id")

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	result, err := h.rejectFriendRequest.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("RejectFriendRequest failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
