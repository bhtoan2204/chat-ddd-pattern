package command

import (
	"context"
	"errors"

	"wechat-clone/core/modules/room/application/dto/in"
	"wechat-clone/core/modules/room/application/dto/out"
	roomsupport "wechat-clone/core/modules/room/application/support"
	roomrepos "wechat-clone/core/modules/room/domain/repos"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"
)

type joinRoomHandler struct {
	baseRepo roomrepos.Repos
}

func NewJoinRoomHandler(baseRepo roomrepos.Repos) cqrs.Handler[*in.JoinRoomRequest, *out.JoinRoomResponse] {
	return &joinRoomHandler{
		baseRepo: baseRepo,
	}
}

func (h *joinRoomHandler) Handle(ctx context.Context, req *in.JoinRoomRequest) (*out.JoinRoomResponse, error) {
	_, err := roomsupport.AccountIDFromCtx(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return nil, stackErr.Error(errors.New("not implemented"))
}
