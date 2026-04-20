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

type unfollowUserHandler struct {
}

func NewUnfollowUser(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
) cqrs.Handler[*in.UnfollowUserRequest, *out.UnfollowUserResponse] {
	return &unfollowUserHandler{}
}

func (u *unfollowUserHandler) Handle(ctx context.Context, req *in.UnfollowUserRequest) (*out.UnfollowUserResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
