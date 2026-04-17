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

type listChatMessagesHandler struct {
	messages roomservice.MessageQueryService
}

func NewListChatMessagesHandler(messages roomservice.MessageQueryService) cqrs.Handler[*in.ListChatMessagesRequest, []*out.ChatMessageResponse] {
	return &listChatMessagesHandler{messages: messages}
}

func (h *listChatMessagesHandler) Handle(ctx context.Context, req *in.ListChatMessagesRequest) ([]*out.ChatMessageResponse, error) {
	accountID, err := roomsupport.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	res, err := h.messages.ListMessages(ctx, accountID, apptypes.ListMessagesQuery{
		RoomID:    req.RoomID,
		Limit:     req.Limit,
		BeforeID:  req.BeforeID,
		BeforeAt:  req.BeforeAt,
		Ascending: req.Ascending,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	outItems := make([]*out.ChatMessageResponse, 0, len(res))
	for _, item := range res {
		copyItem := item
		outItems = append(outItems, roomsupport.ToMessageResponse(&copyItem))
	}

	return outItems, nil
}
