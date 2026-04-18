package command

import (
	"context"
	"reflect"

	"wechat-clone/core/modules/room/application/dto/in"
	"wechat-clone/core/modules/room/application/dto/out"
	"wechat-clone/core/modules/room/application/service"
	roomsupport "wechat-clone/core/modules/room/application/support"
	apptypes "wechat-clone/core/modules/room/application/types"
	roomrepos "wechat-clone/core/modules/room/domain/repos"
	"wechat-clone/core/modules/room/types"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/logging"
	"wechat-clone/core/shared/pkg/stackErr"

	"go.uber.org/zap"
)

type forwardChatMessageHandler struct {
	baseRepo roomrepos.Repos
	realtime service.RealtimeService
}

func NewForwardChatMessageHandler(baseRepo roomrepos.Repos, realtime service.RealtimeService) cqrs.Handler[*in.ForwardChatMessageRequest, *out.ChatMessageResponse] {
	return &forwardChatMessageHandler{baseRepo: baseRepo, realtime: realtime}
}

func (h *forwardChatMessageHandler) Handle(ctx context.Context, req *in.ForwardChatMessageRequest) (*out.ChatMessageResponse, error) {
	accountID, err := roomsupport.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	sourceMessage, err := h.baseRepo.MessageAggregateRepository().Load(ctx, req.MessageID)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	res, err := executeSendMessage(ctx, h.baseRepo, accountID, apptypes.SendMessageCommand{
		RoomID:                 req.TargetRoomID,
		Message:                sourceMessage.Message().Message,
		MessageType:            sourceMessage.Message().MessageType,
		ForwardedFromMessageID: sourceMessage.Message().ID,
		FileName:               sourceMessage.Message().FileName,
		FileSize:               sourceMessage.Message().FileSize,
		MimeType:               sourceMessage.Message().MimeType,
		ObjectKey:              sourceMessage.Message().ObjectKey,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	out := roomsupport.ToMessageResponse(res)
	if err := h.realtime.EmitMessage(ctx, types.MessagePayload{
		RoomId:  out.RoomID,
		Type:    reflect.TypeOf(out).Elem().Name(),
		Payload: out,
	}); err != nil {
		logging.FromContext(ctx).Warnw("failed to emit realtime message after forwarding chat message", zap.Error(err), "message_id", req.MessageID)
	}
	return out, nil
}
