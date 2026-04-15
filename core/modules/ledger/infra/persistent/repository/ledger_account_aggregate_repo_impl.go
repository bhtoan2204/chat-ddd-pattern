package repository

import (
	"context"
	"fmt"
	ledgeraggregate "go-socket/core/modules/ledger/domain/aggregate"
	ledgerrepos "go-socket/core/modules/ledger/domain/repos"
	"go-socket/core/modules/ledger/infra/persistent/model"
	eventpkg "go-socket/core/shared/pkg/event"
	"go-socket/core/shared/pkg/stackErr"
	"reflect"

	"gorm.io/gorm"
)

type ledgerAccountAggregateRepositoryImpl struct {
	db         *gorm.DB
	serializer eventpkg.Serializer
}

func NewLedgerAccountAggregateRepoImpl(db *gorm.DB) ledgerrepos.LedgerAccountAggregateRepository {
	return &ledgerAccountAggregateRepositoryImpl{
		db:         db,
		serializer: newLedgerAccountSerializer(),
	}
}

func newLedgerAccountSerializer() eventpkg.Serializer {
	serializer := eventpkg.NewSerializer()
	if err := serializer.RegisterAggregate(&ledgeraggregate.LedgerAccountAggregate{}); err != nil {
		panic(fmt.Sprintf("register ledger transaction aggregate serializer failed: %v", err))
	}
	return serializer
}

func (r *ledgerAccountAggregateRepositoryImpl) Load(ctx context.Context, accountID string) (*ledgeraggregate.LedgerAccountAggregate, error) {
	if accountID == "" {
		return nil, stackErr.Error(fmt.Errorf("account id is required"))
	}
	aggregate, _ := ledgeraggregate.NewLedgerAccountAggregate(accountID)

	var aggModel model.LedgerAggregateModel
	err := r.db.WithContext(ctx).
		Where("aggregate_id = ? AND aggregate_type = ?", accountID, reflect.TypeOf(aggregate).Elem().Name()).
		First(&aggModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, stackErr.Error(err)
	}

	var eventModels []model.LedgerEventModel
	if err := r.db.WithContext(ctx).
		Where("aggregate_id = ? AND aggregate_type = ?", accountID, reflect.TypeOf(aggregate).Elem().Name()).
		Order("version ASC").
		Find(&eventModels).Error; err != nil {
		return nil, stackErr.Error(err)
	}

	return aggregate, nil
}

func (r *ledgerAccountAggregateRepositoryImpl) Save(ctx context.Context, aggregate *ledgeraggregate.LedgerAccountAggregate) error {
	if aggregate == nil {
		return stackErr.Error(fmt.Errorf("ledger transaction aggregate is nil"))
	}

	root := aggregate.Root()
	events := root.CloneEvents()
	if len(events) == 0 {
		return nil
	}
	if root.BaseVersion() != 0 {
		return stackErr.Error(fmt.Errorf("ledger transaction aggregate does not support updates"))
	}

	root.Update()
	return nil
}
