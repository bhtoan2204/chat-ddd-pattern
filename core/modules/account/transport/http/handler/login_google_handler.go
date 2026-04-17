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

type loginGoogleHandler struct {
	loginGoogle cqrs.Dispatcher[*in.LoginGoogleRequest, *out.LoginGoogleResponse]
}

func NewLoginGoogleHandler(
	loginGoogle cqrs.Dispatcher[*in.LoginGoogleRequest, *out.LoginGoogleResponse],
) *loginGoogleHandler {
	return &loginGoogleHandler{
		loginGoogle: loginGoogle,
	}
}

func (h *loginGoogleHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.LoginGoogleRequest

	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, stackErr.Error(err)
	}

	result, err := h.loginGoogle.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("LoginGoogle failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
