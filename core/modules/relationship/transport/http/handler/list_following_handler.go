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

type listFollowingHandler struct {
	listFollowing cqrs.Dispatcher[*in.ListFollowingRequest, *out.ListFollowingResponse]
}

func NewListFollowingHandler(
	listFollowing cqrs.Dispatcher[*in.ListFollowingRequest, *out.ListFollowingResponse],
) *listFollowingHandler {
	return &listFollowingHandler{
		listFollowing: listFollowing,
	}
}

func (h *listFollowingHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.ListFollowingRequest
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

	result, err := h.listFollowing.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("ListFollowing failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
