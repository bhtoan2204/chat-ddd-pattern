package cqrs

import "context"

type Handler[Req any, Res any] interface {
	Handle(ctx context.Context, req Req) (Res, error)
}

type Dispatcher[Req any, Res any] struct {
	handler Handler[Req, Res]
}

func NewDispatcher[Req any, Res any](handler Handler[Req, Res]) Dispatcher[Req, Res] {
	return Dispatcher[Req, Res]{
		handler: handler,
	}
}

func (d Dispatcher[Req, Res]) Dispatch(ctx context.Context, req Req) (Res, error) {
	return d.handler.Handle(ctx, req)
}
