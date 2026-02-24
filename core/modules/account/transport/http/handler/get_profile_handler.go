// CODE_GENERATOR: handler
package handler

import (
	"errors"
	"go-socket/core/modules/account/application/dto/in"
	"go-socket/core/modules/account/application/query"
	"go-socket/core/shared/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type getProfileHandler struct {
	queryBus query.Bus
}

func NewGetProfileHandler(queryBus query.Bus) *getProfileHandler {
	return &getProfileHandler{
		queryBus: queryBus,
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
	result, err := h.queryBus.GetProfile.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("GetProfile failed", zap.Error(err))
		return nil, errors.New("GetProfile failed")
	}
	return result, nil
}
