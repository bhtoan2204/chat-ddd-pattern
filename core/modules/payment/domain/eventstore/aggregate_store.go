package eventstore

import (
	"context"
	"go-socket/core/shared/pkg/event"

	"gorm.io/gorm"
)

type AggregateStore interface {
	GetAggregate(ctx context.Context, aggregateID string) (event.Aggregate, error)
	SaveAggregate(ctx context.Context, aggregate event.Aggregate) error
}

type aggregateStore struct {
	db *gorm.DB
}
