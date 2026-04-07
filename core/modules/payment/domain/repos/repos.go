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
	ProviderPaymentRepository() ProviderPaymentRepository

	WithTransaction(ctx context.Context, fn func(Repos) error) error
}
