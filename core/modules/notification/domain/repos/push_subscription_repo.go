package repos

import (
	"context"
	"go-socket/core/modules/notification/domain/entity"
)

type PushSubscriptionRepository interface {
	UpsertPushSubscription(ctx context.Context, subscription *entity.PushSubscription) error
	ListPushSubscriptionsByAccountID(ctx context.Context, accountID string) ([]*entity.PushSubscription, error)
}
