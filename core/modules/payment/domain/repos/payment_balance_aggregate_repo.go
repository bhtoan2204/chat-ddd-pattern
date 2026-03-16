package repos

import (
	"context"
	"errors"

	"go-socket/core/modules/payment/domain/aggregate"
)

var ErrPaymentVersionConflict = errors.New("payment aggregate version conflict")

type PaymentBalanceAggregateRepository interface {
	Load(ctx context.Context, accountID string) (*aggregate.PaymentBalanceAggregate, error)
	Save(ctx context.Context, aggregate *aggregate.PaymentBalanceAggregate) error
}
