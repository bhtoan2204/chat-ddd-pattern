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

type listRoomsHandler struct {
	queryBus query.Bus
}

func NewListRoomsHandler(queryBus query.Bus) *listRoomsHandler {
	return &listRoomsHandler{
		queryBus: queryBus,
	}
}

func (h *listRoomsHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.ListRoomsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, err
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, errors.New("validate request failed")
	}
	result, err := h.queryBus.ListRoom.Dispatch(ctx, &request)
	if err != nil {
		logger.Errorw("ListRooms failed", zap.Error(err))
		return nil, errors.New("ListRooms failed")
	}
	return result, nil
}
