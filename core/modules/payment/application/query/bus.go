package query

import (
	"go-socket/core/modules/payment/application/dto/in"
	"go-socket/core/modules/payment/application/dto/out"
	"go-socket/core/shared/pkg/cqrs"
)

type ListTransactionHandler = cqrs.Handler[*in.ListTransactionRequest, *out.ListTransactionResponse]

type Bus struct {
	ListTransaction cqrs.Dispatcher[*in.ListTransactionRequest, *out.ListTransactionResponse]
}

func NewBus(listTransactionHandler ListTransactionHandler) Bus {
	return Bus{
		ListTransaction: cqrs.NewDispatcher(listTransactionHandler),
	}
}
