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

type blockUserHandler struct {
}

func NewBlockUser(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
) cqrs.Handler[*in.BlockUserRequest, *out.BlockUserResponse] {
	return &blockUserHandler{}
}

func (u *blockUserHandler) Handle(ctx context.Context, req *in.BlockUserRequest) (*out.BlockUserResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
