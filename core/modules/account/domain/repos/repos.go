package repos

import "context"

type Repos interface {
	AccountRepository() AccountRepository
	AccountAggregateRepository() AccountAggregateRepository

	WithTransaction(ctx context.Context, fn func(Repos) error) error
}
