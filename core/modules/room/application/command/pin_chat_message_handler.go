package command

import (
	"context"
	"reflect"
	"time"

	"wechat-clone/core/modules/room/application/dto/in"
	"wechat-clone/core/modules/room/application/dto/out"
	"wechat-clone/core/modules/room/application/service"
	roomsupport "wechat-clone/core/modules/room/application/support"
	roomrepos "wechat-clone/core/modules/room/domain/repos"
	"wechat-clone/core/modules/room/types"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

type pinChatMessageHandler struct {
	baseRepo roomrepos.Repos
	realtime service.RealtimeService
}

func NewPinChatMessageHandler(baseRepo roomrepos.Repos, realtime service.RealtimeService) cqrs.Handler[*in.PinChatMessageRequest, *out.ChatConversationResponse] {
	return &pinChatMessageHandler{baseRepo: baseRepo, realtime: realtime}
}

func (h *pinChatMessageHandler) Handle(ctx context.Context, req *in.PinChatMessageRequest) (*out.ChatConversationResponse, error) {
	accountID, err := roomsupport.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	agg, err := h.baseRepo.RoomAggregateRepository().Load(ctx, req.RoomID)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	if err := agg.PinMessage(accountID, req.MessageID, time.Now().UTC(), accountID); err != nil {
		return nil, stackErr.Error(err)
	}
	lastMessage := lastPendingMessage(agg.PendingMessages())
	if err := h.baseRepo.WithTransaction(ctx, func(txRepos roomrepos.Repos) error {
		return stackErr.Error(txRepos.RoomAggregateRepository().Save(ctx, agg))
	}); err != nil {
		return nil, stackErr.Error(err)
	}

	res, err := roomsupport.BuildConversationResultFromState(ctx, h.baseRepo, accountID, agg.Room(), agg.Members(), lastMessage, true)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	out := roomsupport.ToConversationResponse(res)

	if err := h.realtime.EmitMessage(ctx, types.MessagePayload{
		RoomId:  out.RoomID,
		Type:    reflect.TypeOf(out).Elem().Name(),
		Payload: out,
	}); err != nil {
		logging.FromContext(ctx).Warnw("failed to emit realtime message after pinning chat message", zap.Error(err), "message_id", req.MessageID)
	}
	return out, nil
}
