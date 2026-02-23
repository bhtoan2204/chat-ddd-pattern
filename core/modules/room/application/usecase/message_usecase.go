package usecase

import (
	"context"
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
)

type MessageUsecase interface {
	CreateMessage(ctx context.Context, in *in.CreateMessageRequest) (*out.CreateMessageResponse, error)
}
