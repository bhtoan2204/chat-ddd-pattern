// CODE_GENERATOR: handler
package handler

import (
	"errors"
	"go-socket/core/delivery/http/data/in"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/usecase"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type getProfileHandler struct {
	authUsecase usecase.AuthUsecase
}

func NewGetProfileHandler(usecase usecase.Usecase) RequestHandler {
	return &getProfileHandler{
		authUsecase: usecase.AuthUsecase(),
	}
}

func (h *getProfileHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.GetProfileRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, err
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, errors.New("validate request failed")
	}
	result, err := h.authUsecase.GetProfile(ctx, &request)
	if err != nil {
		logger.Errorw("GetProfile failed", zap.Error(err))
		return nil, errors.New("GetProfile failed")
	}
	return result, nil
}
