package command

import (
	"context"
	"fmt"
	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
	"go-socket/core/modules/room/domain/repos"
)

type createMessageHandler struct {
	messageRepo repos.MessageRepository
}

func NewCreateMessageHandler(messageRepo repos.MessageRepository) CreateMessageHandler {
	return &createMessageHandler{
		messageRepo: messageRepo,
	}
}

func (h *createMessageHandler) Handle(ctx context.Context, req *in.CreateMessageRequest) (*out.CreateMessageResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
