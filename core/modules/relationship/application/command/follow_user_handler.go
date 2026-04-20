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

type followUserHandler struct {
}

func NewFollowUser(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
) cqrs.Handler[*in.FollowUserRequest, *out.FollowUserResponse] {
	return &followUserHandler{}
}

func (u *followUserHandler) Handle(ctx context.Context, req *in.FollowUserRequest) (*out.FollowUserResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
