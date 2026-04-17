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

type loginHandler struct {
	login cqrs.Dispatcher[*in.LoginRequest, *out.LoginResponse]
}

func NewLoginHandler(
	login cqrs.Dispatcher[*in.LoginRequest, *out.LoginResponse],
) *loginHandler {
	return &loginHandler{
		login: login,
	}
}

func (h *loginHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.LoginRequest
	request.DeviceUid = c.GetHeader("X-Device-UID")
	request.DeviceName = c.GetHeader("X-Device-Name")
	request.DeviceType = c.GetHeader("X-Device-Type")
	request.OsName = c.GetHeader("X-Device-OS-Name")
	request.OsVersion = c.GetHeader("X-Device-OS-Version")
	request.AppVersion = c.GetHeader("X-Device-App-Version")
	request.UserAgent = c.GetHeader("User-Agent")
	request.IpAddress = c.GetHeader("X-Forwarded-For")
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

	result, err := h.login.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("Login failed", zap.Error(err))
		return nil, stackErr.Error(err)
	}
	return result, nil
}
