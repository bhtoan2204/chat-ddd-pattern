package command

import (
	"context"
	"errors"
	"fmt"
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
	"go-socket/core/modules/room/domain/entity"
	"go-socket/core/modules/room/domain/repos"
	"go-socket/core/shared/infra/xpaseto"
	"go-socket/core/shared/pkg/logging"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type createRoomHandler struct {
	roomRepo repos.RoomRepository
}

func NewCreateRoomHandler(roomRepo repos.RoomRepository) CreateRoomHandler {
	return &createRoomHandler{
		roomRepo: roomRepo,
	}
}

func (h *createRoomHandler) Handle(ctx context.Context, req *in.CreateRoomRequest) (*out.CreateRoomResponse, error) {
	log := logging.FromContext(ctx).Named("CreateRoom")
	account := ctx.Value("account").(*xpaseto.PasetoPayload)
	if account == nil {
		log.Errorw("Account not found", zap.Error(errors.New("account not found")))
		return nil, errors.New("account not found")
	}
	room := &entity.Room{
		ID:          uuid.NewString(),
		Name:        req.Name,
		Description: req.Description,
		RoomType:    req.RoomType,
		OwnerID:     account.AccountID,
	}
	err := h.roomRepo.CreateRoom(ctx, room)
	if err != nil {
		log.Errorw("Failed to create room", zap.Error(err), zap.Any("room", room))
		return nil, fmt.Errorf("create room failed: %w", err)
	}
	return &out.CreateRoomResponse{
		Id:   room.ID,
		Name: room.Name,
	}, nil
}
