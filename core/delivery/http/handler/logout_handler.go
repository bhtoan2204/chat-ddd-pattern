// CODE_GENERATOR: handler
package handler

import (
	"errors"
	"go-socket/core/delivery/http/data/in"
	"go-socket/core/usecase"
	"go-socket/shared/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type logoutHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewLogoutHandler(usecase usecase.Usecase) RequestHandler {
	return &logoutHandler{
		authUsecase: usecase.AuthUsecase(),
	}
}

func (h *logoutHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.LogoutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, err
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, errors.New("validate request failed")
	}
	result, err := h.authUsecase.Logout(ctx, &request)
	if err != nil {
		logger.Errorw("Logout failed", zap.Error(err))
		return nil, errors.New("Logout failed")
	}
	return result, nil
}
