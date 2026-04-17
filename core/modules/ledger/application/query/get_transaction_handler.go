package query

import (
	"context"

	ledgerin "wechat-clone/core/modules/ledger/application/dto/in"
	ledgerout "wechat-clone/core/modules/ledger/application/dto/out"
	ledgerservice "wechat-clone/core/modules/ledger/application/service"
	"wechat-clone/core/shared/pkg/cqrs"
)

type getTransactionHandler struct {
	service ledgerservice.LedgerQueryService
}

func NewGetTransactionHandler(service ledgerservice.LedgerQueryService) cqrs.Handler[*ledgerin.GetTransactionRequest, *ledgerout.TransactionResponse] {
	return &getTransactionHandler{service: service}
}

func (h *getTransactionHandler) Handle(ctx context.Context, req *ledgerin.GetTransactionRequest) (*ledgerout.TransactionResponse, error) {
	return h.service.GetTransaction(ctx, req.TransactionID)
}
