package command

import (
	"context"
	"fmt"
	"time"

	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/ledger/application/dto/in"
	"wechat-clone/core/modules/ledger/application/dto/out"
	ledgerservice "wechat-clone/core/modules/ledger/application/service"
	"wechat-clone/core/shared/finance"
	"wechat-clone/core/shared/infra/lock"
	"wechat-clone/core/shared/pkg/actorctx"
	"wechat-clone/core/shared/pkg/cqrs"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/google/uuid"
)

type transferTransactionHandler struct {
	ledgerService ledgerservice.LedgerService
	locker        lock.Lock
	feePolicy     finance.StripeFeePolicy
	feeAccountID  string
}

func NewTransferTransaction(
	appCtx *appCtx.AppContext,
	ledgerService ledgerservice.LedgerService,
) cqrs.Handler[*in.TransferTransactionRequest, *out.TransactionTransactionResponse] {
	return &transferTransactionHandler{
		ledgerService: ledgerService,
		locker:        appCtx.Locker(),
		feePolicy: finance.StripeFeePolicy{
			RateBPS:    appCtx.GetConfig().LedgerConfig.Stripe.FeeRateBPS,
			FlatAmount: appCtx.GetConfig().LedgerConfig.Stripe.FeeFlatAmount,
		},
		feeAccountID: appCtx.GetConfig().LedgerConfig.Stripe.FeeAccountID,
	}
}

func (u *transferTransactionHandler) Handle(ctx context.Context, req *in.TransferTransactionRequest) (*out.TransactionTransactionResponse, error) {
	fromAccountID, err := actorctx.AccountIDFromContext(ctx)
	if err != nil {
		return nil, stackErr.Error(fmt.Errorf("%w: %w", ledgerservice.ErrUnauthorized, err))
	}

	transferFn := func() (*out.TransactionTransactionResponse, error) {
		feeAmount, err := u.feePolicy.Compute(req.Amount)
		if err != nil {
			return nil, stackErr.Error(fmt.Errorf("%w: %w", ledgerservice.ErrValidation, err))
		}
		transaction, err := u.ledgerService.TransferToAccount(ctx, ledgerservice.TransferToAccountCommand{
			TransactionID: uuid.NewString(),
			FromAccountID: fromAccountID,
			ToAccountID:   req.ToAccountID,
			Currency:      req.Currency,
			Amount:        req.Amount,
			FeeAmount:     feeAmount,
			FeeAccountID:  u.feeAccountID,
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

	lockKeys := []string{fromAccountID, req.ToAccountID}
	if u.feeAccountID != "" {
		lockKeys = append(lockKeys, u.feeAccountID)
	}

	response, err := lock.WithLocks(
		ctx,
		u.locker,
		lockKeys,
		opts,
		transferFn,
	)
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return response, nil
}
