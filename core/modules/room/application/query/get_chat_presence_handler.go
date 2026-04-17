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

type getChatPresenceHandler struct {
	presence roomservice.PresenceQueryService
}

func NewGetChatPresenceHandler(presence roomservice.PresenceQueryService) cqrs.Handler[*in.GetChatPresenceRequest, *out.ChatPresenceResponse] {
	return &getChatPresenceHandler{presence: presence}
}

func (h *getChatPresenceHandler) Handle(ctx context.Context, req *in.GetChatPresenceRequest) (*out.ChatPresenceResponse, error) {
	res, err := h.presence.GetPresence(ctx, apptypes.GetPresenceQuery{AccountID: req.AccountID})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return roomsupport.ToPresenceResponse(res), nil
}
