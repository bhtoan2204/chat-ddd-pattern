package usecase

import (
	"context"
	"errors"
	appCtx "go-socket/core/context"
	"go-socket/core/delivery/http/data/in"
	"go-socket/core/delivery/http/data/out"
	"go-socket/core/domain/room/entity"
	"go-socket/core/domain/room/repos"
	"go-socket/core/shared/infra/xpaseto"
	"time"

	"github.com/google/uuid"
)

type messageUsecaseImpl struct {
	roomRepo    repos.RoomRepository
	messageRepo repos.MessageRepository
}

func NewMessageUsecase(appCtx *appCtx.AppContext, repos repos.Repos) MessageUsecase {
	_ = appCtx
	return &messageUsecaseImpl{
		roomRepo:    repos.RoomRepository(),
		messageRepo: repos.MessageRepository(),
	}
}

func (u *messageUsecaseImpl) CreateMessage(ctx context.Context, in *in.CreateMessageRequest) (*out.CreateMessageResponse, error) {
	if in == nil {
		return nil, errors.New("request is nil")
	}
	if _, err := u.roomRepo.GetRoomByID(ctx, in.RoomID); err != nil {
		return nil, err
	}

	account := ctx.Value("account")
	if account == nil {
		return nil, errors.New("account not found")
	}
	payload, ok := account.(*xpaseto.PasetoPayload)
	if !ok {
		return nil, errors.New("invalid account payload")
	}

	message := &entity.MessageEntity{
		ID:        uuid.New().String(),
		RoomID:    in.RoomID,
		SenderID:  payload.AccountID,
		Message:   in.Message,
		CreatedAt: time.Now().UTC(),
	}
	if err := u.messageRepo.CreateMessage(ctx, message); err != nil {
		return nil, err
	}

	return &out.CreateMessageResponse{
		Id:        message.ID,
		RoomId:    message.RoomID,
		SenderId:  message.SenderID,
		Message:   message.Message,
		CreatedAt: message.CreatedAt.Format(time.RFC3339),
	}, nil
}
