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

type getRoomHandler struct {
	roomUsecase usecase.RoomUsecase
}

func NewGetRoomHandler(usecase usecase.Usecase) RequestHandler {
	return &getRoomHandler{
		roomUsecase: usecase.RoomUsecase(),
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
	result, err := h.roomUsecase.GetRoom(ctx, &request)
	if err != nil {
		logger.Errorw("GetRoom failed", zap.Error(err))
		return nil, errors.New("GetRoom failed")
	}
	return result, nil
}
