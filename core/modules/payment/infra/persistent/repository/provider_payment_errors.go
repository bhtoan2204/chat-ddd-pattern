package repository

import (
	"errors"
	"strings"

	paymentrepos "wechat-clone/core/modules/payment/domain/repos"
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
	if isOracleUniqueConstraintError(err) {
		return stackErr.Error(paymentrepos.ErrProviderPaymentDuplicateIntent)
	}
	return stackErr.Error(err)
}

func isOracleUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToUpper(err.Error()), "ORA-00001")
}
