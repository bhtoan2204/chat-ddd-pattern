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

type unfriendUserHandler struct {
}

func NewUnfriendUser(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
) cqrs.Handler[*in.UnfriendUserRequest, *out.UnfriendUserResponse] {
	return &unfriendUserHandler{}
}

func (u *unfriendUserHandler) Handle(ctx context.Context, req *in.UnfriendUserRequest) (*out.UnfriendUserResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
