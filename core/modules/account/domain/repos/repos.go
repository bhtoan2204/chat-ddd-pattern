package repos

import "context"

type Repos interface {
	AccountRepository() AccountRepository
	AccountOutboxEventsRepository() AccountOutboxEventsRepository

	WithTransaction(ctx context.Context, fn func(Repos) error) error
}
