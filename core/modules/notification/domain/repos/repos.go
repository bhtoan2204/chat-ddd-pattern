package repos

import "context"

type Repos interface {
	NotificationRepository() NotificationRepository

	WithTransaction(ctx context.Context, fn func(Repos) error) error
}
