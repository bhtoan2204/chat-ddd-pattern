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

type acceptFriendRequestHandler struct {
	acceptFriendRequest cqrs.Dispatcher[*in.AcceptFriendRequestRequest, *out.AcceptFriendRequestResponse]
}

func NewAcceptFriendRequestHandler(
	acceptFriendRequest cqrs.Dispatcher[*in.AcceptFriendRequestRequest, *out.AcceptFriendRequestResponse],
) *acceptFriendRequestHandler {
	return &acceptFriendRequestHandler{
		acceptFriendRequest: acceptFriendRequest,
	}
}

func (h *acceptFriendRequestHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.AcceptFriendRequestRequest
	request.RequesterUserID = c.Param("requester_user_id")

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	result, err := h.acceptFriendRequest.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("AcceptFriendRequest failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
