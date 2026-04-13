package repos

import (
	"context"

	"go-socket/core/modules/room/domain/aggregate"
)

type MessageAggregateRepository interface {
	Load(ctx context.Context, messageID string) (*aggregate.MessageStateAggregate, error)
	Save(ctx context.Context, agg *aggregate.MessageStateAggregate) error
}
