package repository

import (
	"errors"

	paymentrepos "wechat-clone/core/modules/payment/domain/repos"
	shareddb "wechat-clone/core/shared/infra/db"
	"wechat-clone/core/shared/pkg/stackErr"

	"gorm.io/gorm"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return stackErr.Error(paymentrepos.ErrProviderPaymentNotFound)
	}
	if shareddb.IsUniqueConstraintError(err) {
		return stackErr.Error(paymentrepos.ErrProviderPaymentDuplicateIntent)
	}
	return stackErr.Error(err)
}
