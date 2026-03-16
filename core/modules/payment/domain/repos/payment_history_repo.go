package repos

import (
	"context"
	"go-socket/core/modules/payment/domain/entity"
	"go-socket/core/shared/utils"
)

type PaymentHistoryRepository interface {
	ListPaymentHistory(ctx context.Context, options utils.QueryOptions) ([]*entity.PaymentHistory, error)
}
