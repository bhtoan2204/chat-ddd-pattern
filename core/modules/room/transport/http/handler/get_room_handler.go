// CODE_GENERATOR: handler
package handler

import (
	"errors"
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/query"
	"go-socket/core/shared/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type getRoomHandler struct {
	queryBus query.Bus
}

func NewGetRoomHandler(queryBus query.Bus) *getRoomHandler {
	return &getRoomHandler{
		queryBus: queryBus,
	}
}

func (h *getRoomHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.GetRoomRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, err
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, errors.New("validate request failed")
	}
	result, err := h.queryBus.GetRoom.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("GetRoom failed", zap.Error(err))
		return nil, errors.New("GetRoom failed")
	}
	return result, nil
}
