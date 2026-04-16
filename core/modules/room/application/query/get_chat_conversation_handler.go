package query

import (
	"context"

	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
	roomservice "go-socket/core/modules/room/application/service"
	roomsupport "go-socket/core/modules/room/application/support"
	apptypes "go-socket/core/modules/room/application/types"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/stackErr"
)

type getChatConversationHandler struct {
	chatService roomservice.Service
}

func NewGetChatConversationHandler(chatService roomservice.Service) cqrs.Handler[*in.GetChatConversationRequest, *out.ChatConversationResponse] {
	return &getChatConversationHandler{chatService: chatService}
}

func (h *getChatConversationHandler) Handle(ctx context.Context, req *in.GetChatConversationRequest) (*out.ChatConversationResponse, error) {
	accountID, err := roomsupport.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	res, err := h.chatService.GetConversation(ctx, accountID, apptypes.GetConversationQuery{
		RoomID: req.RoomID,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return roomsupport.ToConversationResponse(res), nil
}
