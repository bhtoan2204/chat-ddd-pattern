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

type followUserHandler struct {
	followUser cqrs.Dispatcher[*in.FollowUserRequest, *out.FollowUserResponse]
}

func NewFollowUserHandler(
	followUser cqrs.Dispatcher[*in.FollowUserRequest, *out.FollowUserResponse],
) *followUserHandler {
	return &followUserHandler{
		followUser: followUser,
	}
}

func (h *followUserHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.FollowUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	result, err := h.followUser.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("FollowUser failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
