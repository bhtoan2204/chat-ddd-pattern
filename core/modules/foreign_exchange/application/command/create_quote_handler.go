// CODE_GENERATOR: application-handler
package command

import (
	"context"
	"fmt"
	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/foreign_exchange/application/dto/in"
	"wechat-clone/core/modules/foreign_exchange/application/dto/out"
	repos "wechat-clone/core/modules/foreign_exchange/domain/repos"
	"wechat-clone/core/shared/pkg/cqrs"
)

type createQuoteHandler struct {
}

func NewCreateQuote(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
) cqrs.Handler[*in.CreateQuoteRequest, *out.CreateQuoteResponse] {
	return &createQuoteHandler{}
}

func (u *createQuoteHandler) Handle(ctx context.Context, req *in.CreateQuoteRequest) (*out.CreateQuoteResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
