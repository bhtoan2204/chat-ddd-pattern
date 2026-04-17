package query

import (
	"context"

	ledgerin "wechat-clone/core/modules/ledger/application/dto/in"
	ledgerout "wechat-clone/core/modules/ledger/application/dto/out"
	ledgerservice "wechat-clone/core/modules/ledger/application/service"
	"wechat-clone/core/shared/pkg/actorctx"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"
)

type listTransactionHandler struct {
	service ledgerservice.LedgerQueryService
}

func NewListTransactionHandler(service ledgerservice.LedgerQueryService) cqrs.Handler[*ledgerin.ListTransactionRequest, *ledgerout.ListTransactionResponse] {
	return &listTransactionHandler{service: service}
}

func (h *listTransactionHandler) Handle(ctx context.Context, req *ledgerin.ListTransactionRequest) (*ledgerout.ListTransactionResponse, error) {
	accountID, err := actorctx.AccountIDFromContext(ctx)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return h.service.ListTransactions(ctx, accountID, req.Cursor, req.Currency, req.Limit)
}
