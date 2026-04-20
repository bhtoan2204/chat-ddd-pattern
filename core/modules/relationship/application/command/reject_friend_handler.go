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

type rejectFriendRequestHandler struct {
}

func NewRejectFriendRequest(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
) cqrs.Handler[*in.RejectFriendRequestRequest, *out.RejectFriendRequestResponse] {
	return &rejectFriendRequestHandler{}
}

func (u *rejectFriendRequestHandler) Handle(ctx context.Context, req *in.RejectFriendRequestRequest) (*out.RejectFriendRequestResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
