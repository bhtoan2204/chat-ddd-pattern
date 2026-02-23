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

type createRoomHandler struct {
	roomUsecase usecase.RoomUsecase
}

func NewCreateRoomHandler(roomUsecase usecase.RoomUsecase) *createRoomHandler {
	return &createRoomHandler{
		roomUsecase: roomUsecase,
	}
}

func (h *createRoomHandler) Handle(c *gin.Context) (interface{}, error) {
	ctx := c.Request.Context()
	logger := logging.FromContext(ctx)
	var request in.CreateRoomRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Errorw("Unmarshal request failed", zap.Error(err))
		return nil, err
	}
	if err := request.Validate(); err != nil {
		logger.Errorw("Validate request failed", zap.Error(err))
		return nil, errors.New("validate request failed")
	}
	result, err := h.roomUsecase.CreateRoom(ctx, &request)
	if err != nil {
		logger.Errorw("CreateRoom failed", zap.Error(err))
		return nil, errors.New("CreateRoom failed")
	}
	return result, nil
}
