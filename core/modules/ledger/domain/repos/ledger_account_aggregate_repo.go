package repos

import (
	"context"
	"go-socket/core/modules/ledger/domain/aggregate"
)

type LedgerAccountAggregateRepository interface {
	Load(ctx context.Context, accountID string) (*aggregate.LedgerAccountAggregate, error)
	Save(ctx context.Context, aggregate *aggregate.LedgerAccountAggregate) error
}
