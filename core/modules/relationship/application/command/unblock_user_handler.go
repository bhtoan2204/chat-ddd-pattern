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

type unblockUserHandler struct {
}

func NewUnblockUser(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
) cqrs.Handler[*in.UnblockUserRequest, *out.UnblockUserResponse] {
	return &unblockUserHandler{}
}

func (u *unblockUserHandler) Handle(ctx context.Context, req *in.UnblockUserRequest) (*out.UnblockUserResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
