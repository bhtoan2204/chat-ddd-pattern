package usecase

import (
	"context"
	"errors"
	"fmt"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
	"go-socket/core/modules/room/domain/entity"
	"go-socket/core/modules/room/domain/repos"
	"go-socket/core/modules/room/types"
	"go-socket/core/shared/infra/xpaseto"
	"go-socket/core/shared/pkg/logging"
	"go-socket/utils"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type roomUsecaseImpl struct {
	roomRepo repos.RoomRepository
}

func NewRoomUsecase(appCtx *appCtx.AppContext, repos repos.Repos) RoomUsecase {
	_ = appCtx
	return &roomUsecaseImpl{roomRepo: repos.RoomRepository()}
}

func (u *roomUsecaseImpl) CreateRoom(ctx context.Context, in *in.CreateRoomRequest) (*out.CreateRoomResponse, error) {
	log := logging.FromContext(ctx).Named("CreateRoom")
	account := ctx.Value("account").(*xpaseto.PasetoPayload)
	if account == nil {
		log.Errorw("Account not found", zap.Error(errors.New("account not found")))
		return nil, errors.New("account not found")
	}
	room := &entity.Room{
		ID:          uuid.NewString(),
		Name:        in.Name,
		Description: in.Description,
		RoomType:    types.RoomType(in.RoomType),
		OwnerID:     account.AccountID,
	}
	err := u.roomRepo.CreateRoom(ctx, room)
	if err != nil {
		log.Errorw("Failed to create room", zap.Error(err))
		return nil, fmt.Errorf("create room failed: %w", err)
	}
	return &out.CreateRoomResponse{
		Id:   room.ID,
		Name: room.Name,
	}, nil
}

func (u *roomUsecaseImpl) ListRooms(ctx context.Context, in *in.ListRoomsRequest) (*out.ListRoomsResponse, error) {
	rooms, err := u.roomRepo.ListRooms(ctx, utils.QueryOptions{
		Conditions:     []utils.Condition{},
		Limit:          lo.ToPtr(in.Limit),
		Offset:         lo.ToPtr((in.Page - 1) * in.Limit),
		OrderBy:        "created_at",
		OrderDirection: "desc",
	})
	if err != nil {
		return nil, err
	}
	return &out.ListRoomsResponse{
		Rooms: lo.Map(rooms, func(room *entity.Room, _ int) out.RoomResponse {
			return out.RoomResponse{
				Id:          room.ID,
				Name:        room.Name,
				Description: room.Description,
				RoomType:    string(room.RoomType),
				OwnerId:     room.OwnerID,
				CreatedAt:   room.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   room.UpdatedAt.Format(time.RFC3339),
			}
		}),
	}, nil
}

func (u *roomUsecaseImpl) GetRoom(ctx context.Context, in *in.GetRoomRequest) (*out.GetRoomResponse, error) {
	room, err := u.roomRepo.GetRoomByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &out.GetRoomResponse{
		Id:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		RoomType:    string(room.RoomType),
		OwnerId:     room.OwnerID,
		CreatedAt:   room.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   room.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *roomUsecaseImpl) UpdateRoom(ctx context.Context, in *in.UpdateRoomRequest) (*out.UpdateRoomResponse, error) {
	room, err := u.roomRepo.GetRoomByID(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	room.Name = in.Name
	err = u.roomRepo.UpdateRoom(ctx, room)
	if err != nil {
		return nil, err
	}
	return &out.UpdateRoomResponse{
		Id:        room.ID,
		Name:      room.Name,
		CreatedAt: room.CreatedAt.Format(time.RFC3339),
		UpdatedAt: room.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (u *roomUsecaseImpl) DeleteRoom(ctx context.Context, in *in.DeleteRoomRequest) (*out.DeleteRoomResponse, error) {
	err := u.roomRepo.DeleteRoom(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &out.DeleteRoomResponse{
		Message: "Room deleted successfully",
	}, nil
}
