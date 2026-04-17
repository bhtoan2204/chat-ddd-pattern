// CODE_GENERATOR - do not edit: handler
package handler

import (
	"net/http"

	"wechat-clone/core/modules/account/application/dto/in"
	"wechat-clone/core/modules/account/application/dto/out"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type getAvatarHandler struct {
	getAvatar cqrs.Dispatcher[*in.GetAvatarRequest, *out.GetAvatarResponse]
}

func NewGetAvatarHandler(
	getAvatar cqrs.Dispatcher[*in.GetAvatarRequest, *out.GetAvatarResponse],
) *getAvatarHandler {
	return &getAvatarHandler{
		getAvatar: getAvatar,
	}
}

func (h *getAvatarHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.GetAvatarRequest
	request.AccountID = c.Param("account_id")

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	result, err := h.getAvatar.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("GetAvatar failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
