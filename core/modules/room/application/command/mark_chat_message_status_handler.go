package command

import (
	"context"
	"errors"
	"reflect"
	"time"

	"wechat-clone/core/modules/room/application/dto/in"
	"wechat-clone/core/modules/room/application/dto/out"
	"wechat-clone/core/modules/room/application/service"
	roomsupport "wechat-clone/core/modules/room/application/support"
	"wechat-clone/core/modules/room/domain/entity"
	roomrepos "wechat-clone/core/modules/room/domain/repos"
	"wechat-clone/core/modules/room/types"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type markChatMessageStatusHandler struct {
	baseRepo roomrepos.Repos
	realtime service.RealtimeService
}

func NewMarkChatMessageStatusHandler(baseRepo roomrepos.Repos, realtime service.RealtimeService) cqrs.Handler[*in.MarkChatMessageStatusRequest, *out.MarkChatMessageStatusResponse] {
	return &markChatMessageStatusHandler{baseRepo: baseRepo, realtime: realtime}
}

func (h *markChatMessageStatusHandler) Handle(ctx context.Context, req *in.MarkChatMessageStatusRequest) (*out.MarkChatMessageStatusResponse, error) {
	accountID, err := roomsupport.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	agg, err := h.baseRepo.MessageAggregateRepository().Load(ctx, req.MessageID)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	var member = (*entity.RoomMemberEntity)(nil)
	if roomMember, memberErr := h.baseRepo.RoomMemberRepository().GetRoomMemberByAccount(ctx, agg.Message().RoomID, accountID); memberErr == nil {
		member = roomMember
	} else if !errors.Is(memberErr, gorm.ErrRecordNotFound) {
		return nil, stackErr.Error(memberErr)
	}

	changed, err := agg.MarkStatus(accountID, req.Status, member, time.Now().UTC())
	if err != nil {
		return nil, stackErr.Error(err)
	}
	if changed {
		if err := h.baseRepo.WithTransaction(ctx, func(txRepos roomrepos.Repos) error {
			return stackErr.Error(txRepos.MessageAggregateRepository().Save(ctx, agg))
		}); err != nil {
			return nil, stackErr.Error(err)
		}
	}

	out := &out.MarkChatMessageStatusResponse{Ok: true}
	if err := h.realtime.EmitMessage(ctx, types.MessagePayload{
		RoomId:  agg.Message().RoomID,
		Type:    reflect.TypeOf(out).Elem().Name(),
		Payload: out,
	}); err != nil {
		logging.FromContext(ctx).Warnw("failed to emit realtime message after marking chat message status", zap.Error(err), "message_id", req.MessageID)
	}
	return out, nil
}
