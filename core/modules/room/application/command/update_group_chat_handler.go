package command

import (
	"context"
	"reflect"
	"time"

	"wechat-clone/core/modules/room/application/dto/in"
	"wechat-clone/core/modules/room/application/dto/out"
	"wechat-clone/core/modules/room/application/service"
	roomsupport "wechat-clone/core/modules/room/application/support"
	"wechat-clone/core/modules/room/domain/aggregate"
	roomrepos "wechat-clone/core/modules/room/domain/repos"
	"wechat-clone/core/modules/room/types"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"
)

type updateGroupChatHandler struct {
	baseRepo roomrepos.Repos
	realtime service.RealtimeService
}

func NewUpdateGroupChatHandler(baseRepo roomrepos.Repos, realtime service.RealtimeService) cqrs.Handler[*in.UpdateGroupChatRequest, *out.ChatConversationResponse] {
	return &updateGroupChatHandler{baseRepo: baseRepo, realtime: realtime}
}
func (h *updateGroupChatHandler) Handle(ctx context.Context, req *in.UpdateGroupChatRequest) (*out.ChatConversationResponse, error) {
	accountID, err := roomsupport.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	agg, err := h.baseRepo.RoomAggregateRepository().Load(ctx, req.RoomID)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	updated, err := agg.UpdateGroupDetails(aggregate.UpdateGroupDetailsParams{
		ActorID:       accountID,
		Name:          req.Name,
		Description:   req.Description,
		Now:           time.Now().UTC(),
		SystemActorID: accountID,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}
	lastMessage := lastPendingMessage(agg.PendingMessages())
	if updated {
		if err := h.baseRepo.WithTransaction(ctx, func(txRepos roomrepos.Repos) error {
			return stackErr.Error(txRepos.RoomAggregateRepository().Save(ctx, agg))
		}); err != nil {
			return nil, stackErr.Error(err)
		}
	}

	res, err := roomsupport.BuildConversationResultFromState(ctx, h.baseRepo, accountID, agg.Room(), agg.Members(), lastMessage, true)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	out := roomsupport.ToConversationResponse(res)

	h.realtime.EmitMessage(ctx, types.MessagePayload{
		RoomId:  out.RoomID,
		Type:    reflect.TypeOf(out).Elem().Name(),
		Payload: out,
	})

	return out, nil
}
