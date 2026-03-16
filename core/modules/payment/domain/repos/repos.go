package repos

import (
	"context"
)

type Repos interface {
	PaymentBalanceAggregateRepository() PaymentBalanceAggregateRepository
	PaymentProjectionRepository() PaymentProjectionRepository
	PaymentOutboxEventsRepository() PaymentOutboxEventsRepository
	PaymentAccountProjectionRepository() PaymentAccountProjectionRepository
	PaymentHistoryRepository() PaymentHistoryRepository

	WithTransaction(ctx context.Context, fn func(Repos) error) error
}
