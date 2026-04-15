package command

import (
	"context"
	"fmt"
	appCtx "go-socket/core/context"
	"go-socket/core/modules/ledger/application/dto/in"
	"go-socket/core/modules/ledger/application/dto/out"
	"go-socket/core/modules/ledger/application/service"
	ledgeraggregate "go-socket/core/modules/ledger/domain/aggregate"
	repos "go-socket/core/modules/ledger/domain/repos"
	"go-socket/core/shared/pkg/cqrs"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type transferTransactionHandler struct {
}

func NewTransferTransaction(
	appCtx *appCtx.AppContext,
	baseRepo repos.Repos,
	service *service.LedgerService,
) cqrs.Handler[*in.TransferTransactionRequest, *out.TransactionTransactionResponse] {
	return &transferTransactionHandler{}
}

func (u *transferTransactionHandler) Handle(ctx context.Context, req *in.TransferTransactionRequest) (*out.TransactionTransactionResponse, error) {
	log := logging.FromContext(ctx)
	transactionID := uuid.New().String()
	aggregate, err := ledgeraggregate.NewLedgerTransactionAggregate(transactionID)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	log.Infow("data", zap.Any("aggregate", aggregate))
	return nil, fmt.Errorf("not implemented yet")
}
