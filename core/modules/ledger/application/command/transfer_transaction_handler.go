package command

import (
	"context"
	"fmt"
	"time"

	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/ledger/application/dto/in"
	"wechat-clone/core/modules/ledger/application/dto/out"
	ledgerservice "wechat-clone/core/modules/ledger/application/service"
	"wechat-clone/core/shared/infra/lock"
	"wechat-clone/core/shared/pkg/actorctx"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/google/uuid"
)

type transferTransactionHandler struct {
	ledgerService ledgerservice.LedgerService
	locker        lock.Lock
}

func NewTransferTransaction(
	appCtx *appCtx.AppContext,
	ledgerService ledgerservice.LedgerService,
) cqrs.Handler[*in.TransferTransactionRequest, *out.TransactionTransactionResponse] {
	return &transferTransactionHandler{
		ledgerService: ledgerService,
		locker:        appCtx.Locker(),
	}
}

func (u *transferTransactionHandler) Handle(ctx context.Context, req *in.TransferTransactionRequest) (*out.TransactionTransactionResponse, error) {
	fromAccountID, err := actorctx.AccountIDFromContext(ctx)
	if err != nil {
		return nil, stackErr.Error(fmt.Errorf("%w: %w", ledgerservice.ErrUnauthorized, err))
	}

	transferFn := func() (*out.TransactionTransactionResponse, error) {
		transaction, err := u.ledgerService.TransferToAccount(ctx, ledgerservice.TransferToAccountCommand{
			TransactionID: uuid.NewString(),
			FromAccountID: fromAccountID,
			ToAccountID:   req.ToAccountID,
			Currency:      req.Currency,
			Amount:        req.Amount,
		})
		if err != nil {
			return nil, stackErr.Error(err)
		}

		responseEntries := make([]out.LedgerEntryResponse, 0, len(transaction.Entries))
		for _, entry := range transaction.Entries {
			responseEntries = append(responseEntries, out.LedgerEntryResponse{
				ID:            entry.ID,
				TransactionID: entry.TransactionID,
				AccountID:     entry.AccountID,
				Currency:      entry.Currency,
				Amount:        entry.Amount,
				CreatedAt:     entry.CreatedAt,
			})
		}

		return &out.TransactionTransactionResponse{
			TransactionID: transaction.TransactionID,
			Currency:      transaction.Currency,
			CreatedAt:     transaction.CreatedAt.UTC().Format(time.RFC3339Nano),
			Entries:       responseEntries,
		}, nil
	}

	opts := lock.DefaultMultiLockOptions()
	opts.KeyPrefix = ledgerservice.LedgerAccountLockKeyPrefix

	response, err := lock.WithLocks(
		ctx,
		u.locker,
		[]string{fromAccountID, req.ToAccountID},
		opts,
		transferFn,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return response, nil
}
