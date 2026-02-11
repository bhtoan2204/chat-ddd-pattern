package usecase

import (
	"context"
	"go-socket/core/delivery/http/data/in"
	"go-socket/core/delivery/http/data/out"
)

type MessageUsecase interface {
	CreateMessage(ctx context.Context, in *in.CreateMessageRequest) (*out.CreateMessageResponse, error)
}
