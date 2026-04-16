// CODE_GENERATOR: application-handler
package command

import (
	"context"
	"fmt"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/ledger/application/dto/in"
	"go-socket/core/modules/ledger/application/dto/out"
	"go-socket/core/modules/ledger/application/service"
	repos "go-socket/core/modules/ledger/domain/repos"
	"go-socket/core/shared/pkg/cqrs"
)

type createTransactionHandler struct {
}

func NewCreateTransaction(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
	service *service.LedgerService,
) cqrs.Handler[*in.CreateTransactionRequest, *out.TransactionResponse] {
	return &createTransactionHandler{}
}

func (u *createTransactionHandler) Handle(ctx context.Context, req *in.CreateTransactionRequest) (*out.TransactionResponse, error) {
	return nil, fmt.Errorf("not implemented yet")
}
