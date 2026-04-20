package repository

import (
	"errors"

	ledgerrepos "wechat-clone/core/modules/ledger/domain/repos"
	shareddb "wechat-clone/core/shared/infra/db"
	"wechat-clone/core/shared/pkg/stackErr"

	"gorm.io/gorm"
)

var (
	ErrNotFound  = ledgerrepos.ErrNotFound
	ErrDuplicate = ledgerrepos.ErrDuplicate
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return stackErr.Error(ErrNotFound)
	}

	if shareddb.IsUniqueConstraintError(err) {
		return stackErr.Error(ErrDuplicate)
	}
	return stackErr.Error(err)
}
