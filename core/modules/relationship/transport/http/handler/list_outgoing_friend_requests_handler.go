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

type listOutgoingFriendRequestsHandler struct {
	listOutgoingFriendRequests cqrs.Dispatcher[*in.ListOutgoingFriendRequestsRequest, *out.ListOutgoingFriendRequestsResponse]
}

func NewListOutgoingFriendRequestsHandler(
	listOutgoingFriendRequests cqrs.Dispatcher[*in.ListOutgoingFriendRequestsRequest, *out.ListOutgoingFriendRequestsResponse],
) *listOutgoingFriendRequestsHandler {
	return &listOutgoingFriendRequestsHandler{
		listOutgoingFriendRequests: listOutgoingFriendRequests,
	}
}

func (h *listOutgoingFriendRequestsHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.ListOutgoingFriendRequestsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	result, err := h.listOutgoingFriendRequests.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("ListOutgoingFriendRequests failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
