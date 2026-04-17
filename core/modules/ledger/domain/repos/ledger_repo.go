package repos

import (
	"context"
	"time"

	"wechat-clone/core/modules/ledger/domain/entity"
)

type ListTransactionsFilter struct {
	AccountID           string
	Currency            string
	CursorCreatedAt     *time.Time
	CursorTransactionID string
	Limit               int
}

// LedgerRepository exposes read-side ledger views derived from canonical
// transaction postings. Write-side persistence must go through aggregate
// repositories to keep the posting model explicit.
//
//go:generate mockgen -package=repos -destination=ledger_repo_mock.go -source=ledger_repo.go
type LedgerRepository interface {
	GetBalance(ctx context.Context, accountID, currency string) (int64, error)
	GetTransaction(ctx context.Context, transactionID string) (*entity.LedgerTransaction, error)
	ListTransactions(ctx context.Context, filter ListTransactionsFilter) ([]*entity.LedgerTransaction, error)
	CountTransactions(ctx context.Context, accountID, currency string) (int64, error)
}
