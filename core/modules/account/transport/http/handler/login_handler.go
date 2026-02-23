// CODE_GENERATOR: handler
package handler

import (
	"errors"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/usecase"
	"go-socket/core/shared/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type loginHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewLoginHandler(authUsecase usecase.AuthUsecase) *loginHandler {
	return &loginHandler{
		authUsecase: authUsecase,
	}
}

func (h *loginHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, err
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, errors.New("validate request failed")
	}
	result, err := h.authUsecase.Login(ctx, &request)
	if err != nil {
		logger.Errorw("Login failed", zap.Error(err))
		return nil, errors.New("Login failed")
	}
	return result, nil
}
