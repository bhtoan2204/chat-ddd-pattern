package repos

import (
	"context"
	"wechat-clone/core/modules/relationship/domain/aggregate"
)

type FriendRequestAggregateRepository interface {
	Load(ctx context.Context, friendRequestID string) (*aggregate.FriendRequestAggregate, error)
	Save(ctx context.Context, agg *aggregate.FriendRequestAggregate) error
}
