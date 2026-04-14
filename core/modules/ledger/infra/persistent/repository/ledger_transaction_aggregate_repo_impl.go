package repository

import (
	"context"
	"fmt"

	ledgeraggregate "go-socket/core/modules/ledger/domain/aggregate"
	"go-socket/core/modules/ledger/domain/entity"
	ledgerrepos "go-socket/core/modules/ledger/domain/repos"
	"go-socket/core/modules/ledger/infra/persistent/model"
	"go-socket/core/shared/pkg/stackErr"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type ledgerTransactionAggregateRepoImpl struct {
	db *gorm.DB
}

func NewLedgerTransactionAggregateRepoImpl(db *gorm.DB) ledgerrepos.LedgerTransactionAggregateRepository {
	return &ledgerTransactionAggregateRepoImpl{db: db}
}

func (r *ledgerTransactionAggregateRepoImpl) Save(ctx context.Context, aggregate *ledgeraggregate.LedgerTransactionAggregate) error {
	if aggregate == nil {
		return stackErr.Error(fmt.Errorf("ledger transaction aggregate is nil"))
	}

	transaction, err := aggregate.Snapshot()
	if err != nil {
		return stackErr.Error(err)
	}

	if err := r.db.WithContext(ctx).Create(&model.LedgerTransactionModel{
		TransactionID: transaction.TransactionID,
		CreatedAt:     transaction.CreatedAt,
	}).Error; err != nil {
		return mapError(err)
	}

	entryModels := lo.Map(transaction.Entries, func(entry *entity.LedgerEntry, _ int) model.LedgerEntryModel {
		return model.LedgerEntryModel{
			TransactionID: entry.TransactionID,
			AccountID:     entry.AccountID,
			Amount:        entry.Amount,
			CreatedAt:     entry.CreatedAt,
		}
	})

	if err := r.db.WithContext(ctx).Create(&entryModels).Error; err != nil {
		return mapError(err)
	}

	if err := aggregate.AssignEntryIDs(lo.Map(entryModels, func(entry model.LedgerEntryModel, _ int) int64 {
		return entry.ID
	})); err != nil {
		return stackErr.Error(err)
	}

	return nil
}
