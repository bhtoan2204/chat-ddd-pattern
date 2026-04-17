package query

import (
	"context"

	"wechat-clone/core/modules/room/application/dto/in"
	"wechat-clone/core/modules/room/application/dto/out"
	roomservice "wechat-clone/core/modules/room/application/service"
	roomsupport "wechat-clone/core/modules/room/application/support"
	apptypes "wechat-clone/core/modules/room/application/types"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"
)

type getChatConversationHandler struct {
	conversations roomservice.ConversationQueryService
}

func NewGetChatConversationHandler(conversations roomservice.ConversationQueryService) cqrs.Handler[*in.GetChatConversationRequest, *out.ChatConversationResponse] {
	return &getChatConversationHandler{conversations: conversations}
}

func (h *getChatConversationHandler) Handle(ctx context.Context, req *in.GetChatConversationRequest) (*out.ChatConversationResponse, error) {
	accountID, err := roomsupport.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	res, err := h.conversations.GetConversation(ctx, accountID, apptypes.GetConversationQuery{
		RoomID: req.RoomID,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return roomsupport.ToConversationResponse(res), nil
}
