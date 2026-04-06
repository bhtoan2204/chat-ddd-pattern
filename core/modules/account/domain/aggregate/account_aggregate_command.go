package aggregate

import (
	"reflect"

	"go-socket/core/shared/pkg/stackErr"
)

func NewAccountAggregate(accountID string) (*AccountAggregate, error) {
	agg := &AccountAggregate{}
	agg.SetAggregateType(reflect.TypeOf(agg).Elem().Name())
	if err := agg.SetID(accountID); err != nil {
		return nil, stackErr.Error(err)
	}

	return agg, nil
}
