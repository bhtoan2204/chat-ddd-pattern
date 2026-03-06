package repos

import "context"

type Repos interface {
	NotificationRepository() NotificationRepository
	PushSubscriptionRepository() PushSubscriptionRepository

	WithTransaction(ctx context.Context, fn func(Repos) error) error
}
