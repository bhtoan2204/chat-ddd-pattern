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

type sendFriendRequestHandler struct {
}

func NewSendFriendRequest(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
) cqrs.Handler[*in.SendFriendRequestRequest, *out.SendFriendRequestResponse] {
	return &sendFriendRequestHandler{}
}

func (u *sendFriendRequestHandler) Handle(ctx context.Context, req *in.SendFriendRequestRequest) (*out.SendFriendRequestResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
