package repos

import (
	"context"

	"go-socket/core/modules/ledger/domain/aggregate"
)

//go:generate mockgen -package=repos -destination=ledger_transaction_aggregate_repo_mock.go -source=ledger_transaction_aggregate_repo.go
type LedgerTransactionAggregateRepository interface {
	Save(ctx context.Context, aggregate *aggregate.LedgerTransactionAggregate) error
}
