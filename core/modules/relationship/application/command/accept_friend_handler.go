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

type acceptFriendRequestHandler struct {
}

func NewAcceptFriendRequest(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
) cqrs.Handler[*in.AcceptFriendRequestRequest, *out.AcceptFriendRequestResponse] {
	return &acceptFriendRequestHandler{}
}

func (u *acceptFriendRequestHandler) Handle(ctx context.Context, req *in.AcceptFriendRequestRequest) (*out.AcceptFriendRequestResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
