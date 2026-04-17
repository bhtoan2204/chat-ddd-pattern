package repository

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	ledgeraggregate "go-socket/core/modules/ledger/domain/aggregate"
	ledgerrepos "go-socket/core/modules/ledger/domain/repos"
	eventpkg "go-socket/core/shared/pkg/event"
	"go-socket/core/shared/pkg/stackErr"
)

type aggregateStore interface {
	Get(ctx context.Context, aggregateID string, agg eventpkg.Aggregate) error
	FindPostedTransaction(ctx context.Context, aggregateID, aggregateType, transactionID string) (*ledgeraggregate.LedgerAccountPosting, error)
	Save(ctx context.Context, agg eventpkg.Aggregate) error
}

type aggregateStoreImpl struct {
	repo ledgerEventStore
}

func newAggregateStore(dbTX dbTX, serializer eventpkg.Serializer) aggregateStore {
	return &aggregateStoreImpl{
		repo: newLedgerEventStore(dbTX, serializer),
	}
}

func (s *aggregateStoreImpl) Get(ctx context.Context, aggregateID string, agg eventpkg.Aggregate) error {
	if reflect.ValueOf(agg).Kind() != reflect.Ptr {
		return stackErr.Error(fmt.Errorf("aggregate must be a pointer"))
	}

	aggregateType := eventpkg.AggregateTypeName(agg)
	agg.Root().SetAggregateType(aggregateType)

	hasSnapshot, err := s.repo.ReadSnapshot(ctx, aggregateID, aggregateType, agg)
	if err != nil {
		return stackErr.Error(err)
	}
	if !hasSnapshot {
		return stackErr.Error(s.repo.Get(ctx, aggregateID, aggregateType, 0, agg))
	}

	return stackErr.Error(s.repo.Get(ctx, aggregateID, aggregateType, agg.Root().BaseVersion(), agg))
}

func (s *aggregateStoreImpl) FindPostedTransaction(
	ctx context.Context,
	aggregateID string,
	aggregateType string,
	transactionID string,
) (*ledgeraggregate.LedgerAccountPosting, error) {
	return s.repo.FindPostedTransaction(ctx, aggregateID, aggregateType, transactionID)
}

func (s *aggregateStoreImpl) Save(ctx context.Context, agg eventpkg.Aggregate) error {
	if reflect.ValueOf(agg).Kind() != reflect.Ptr {
		return stackErr.Error(fmt.Errorf("aggregate must be a pointer"))
	}

	root := agg.Root()
	aggregateType := eventpkg.AggregateTypeName(agg)
	root.SetAggregateType(aggregateType)

	events := root.CloneEvents()
	if len(events) == 0 {
		return nil
	}

	if err := s.repo.CreateIfNotExist(ctx, root.AggregateID(), aggregateType); err != nil {
		return stackErr.Error(err)
	}

	for idx, evt := range events {
		if err := s.repo.ReservePostedTransaction(ctx, evt); err != nil {
			if errors.Is(err, ledgerrepos.ErrAlreadyApplied) {
				if len(events) == 1 && idx == 0 {
					return stackErr.Error(err)
				}
				return stackErr.Error(fmt.Errorf("ledger idempotency collision on event #%d: %w", idx, err))
			}
			return stackErr.Error(fmt.Errorf("reserve ledger posting #%d failed: %w", idx, err))
		}
	}

	if ok, err := s.repo.CheckAndUpdateVersion(ctx, root.AggregateID(), aggregateType, root.BaseVersion(), root.Version()); err != nil {
		return stackErr.Error(err)
	} else if !ok {
		return stackErr.Error(fmt.Errorf(
			"optimistic concurrency control failed id=%s expectedVersion=%d newVersion=%d",
			root.AggregateID(),
			root.BaseVersion(),
			root.Version(),
		))
	}

	for idx, evt := range events {
		if err := s.repo.Append(ctx, evt); err != nil {
			return stackErr.Error(fmt.Errorf("append ledger event #%d failed: %w", idx, err))
		}
		if evt.Version%100 == 0 {
			if err := s.repo.CreateSnapshot(ctx, agg); err != nil {
				return stackErr.Error(fmt.Errorf("create ledger snapshot failed: %w", err))
			}
		}
	}

	root.Update()
	return nil
}
