// CODE_GENERATOR: application-handler
package command

import (
	"context"
	"fmt"
	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/relationship/application/dto/in"
	"wechat-clone/core/modules/relationship/application/dto/out"
	repos "wechat-clone/core/modules/relationship/domain/repos"
	"wechat-clone/core/shared/pkg/cqrs"
)

type cancelFriendRequestHandler struct {
}

func NewCancelFriendRequest(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
) cqrs.Handler[*in.CancelFriendRequestRequest, *out.CancelFriendRequestResponse] {
	return &cancelFriendRequestHandler{}
}

func (u *cancelFriendRequestHandler) Handle(ctx context.Context, req *in.CancelFriendRequestRequest) (*out.CancelFriendRequestResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
