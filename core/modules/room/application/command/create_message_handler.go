package command

import (
	"context"

	"go-socket/core/modules/room/application/dto/in"
	"go-socket/core/modules/room/application/dto/out"
	roomsupport "go-socket/core/modules/room/application/support"
	apptypes "go-socket/core/modules/room/application/types"
	roomrepos "go-socket/core/modules/room/domain/repos"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/stackErr"
)

type createMessageHandler struct {
	baseRepo roomrepos.Repos
}

func NewCreateMessageHandler(baseRepo roomrepos.Repos) cqrs.Handler[*in.CreateMessageRequest, *out.CreateMessageResponse] {
	return &createMessageHandler{baseRepo: baseRepo}
}

func (h *createMessageHandler) Handle(ctx context.Context, req *in.CreateMessageRequest) (*out.CreateMessageResponse, error) {
	accountID, err := roomsupport.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	res, err := executeSendMessage(ctx, h.baseRepo, accountID, apptypes.SendMessageCommand{
		RoomID:  req.RoomID,
		Message: req.Message,
	})
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return &out.CreateMessageResponse{
		ID:        res.ID,
		RoomID:    res.RoomID,
		SenderID:  res.SenderID,
		Message:   res.Message,
		CreatedAt: res.CreatedAt,
	}, nil
}
