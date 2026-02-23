// CODE_GENERATOR: handler
package handler

import (
	"errors"
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/usecase"
	"go-socket/core/shared/pkg/logging"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type updateRoomHandler struct {
	roomUsecase usecase.RoomUsecase
}

func NewUpdateRoomHandler(roomUsecase usecase.RoomUsecase) *updateRoomHandler {
	return &updateRoomHandler{
		roomUsecase: roomUsecase,
	}
}

func (h *updateRoomHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.UpdateRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, err
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, errors.New("validate request failed")
	}
	result, err := h.roomUsecase.UpdateRoom(ctx, &request)
	if err != nil {
		logger.Errorw("UpdateRoom failed", zap.Error(err))
		return nil, errors.New("UpdateRoom failed")
	}
	return result, nil
}
