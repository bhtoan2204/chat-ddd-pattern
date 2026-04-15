package aggregate

import (
	"go-socket/core/modules/ledger/domain/entity"
	"go-socket/core/shared/pkg/event"
	"go-socket/core/shared/pkg/stackErr"
	"reflect"
	"strings"
)

type LedgerAccountAggregate struct {
	event.AggregateRoot

	Entries []*entity.LedgerEntry
}

// RegisterEvents implements [event.BaseAggregate].
func (l *LedgerAccountAggregate) RegisterEvents(event.RegisterEventsFunc) error {
	panic("unimplemented")
}

func NewLedgerAccountAggregate(accountID string) (*LedgerAccountAggregate, error) {
	agg := &LedgerAccountAggregate{}
	agg.Root().SetAggregateType(reflect.TypeOf(agg).Elem().Name())
	if err := agg.SetID(strings.TrimSpace(accountID)); err != nil {
		return nil, stackErr.Error(err)
	}

	return agg, nil
}

// func (a *LedgerAccountAggregate) Load
