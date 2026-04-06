package repos

import (
	"context"
	"go-socket/core/modules/account/domain/aggregate"
)

type AccountAggregateRepository interface {
	Load(ctx context.Context, accountID string) (*aggregate.AccountAggregate, error)
	Save(ctx context.Context, agg *aggregate.AccountAggregate) error
}
