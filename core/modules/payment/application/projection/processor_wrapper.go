package projection

import (
	"context"
	"fmt"

	infraMessaging "go-socket/core/shared/infra/messaging"
	"go-socket/core/shared/pkg/contxt"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"
)

func (p *processor) processMessage(consume infraMessaging.Consumer) infraMessaging.CallBack {
	return func(ctx context.Context, _ string, vals []byte) (err error) {
		ctx = contxt.SetRequestID(ctx)

		logger := logging.FromContext(ctx)
		if reqID := contxt.RequestIDFromCtx(ctx); reqID != "" {
			logger = logger.With("request_id", reqID)
		}
		ctx = logging.WithLogger(ctx, logger)

		defer func() {
			if r := recover(); r != nil {
				err = stackErr.Error(fmt.Errorf("panic recovered: %v", r))
			}
		}()

		handler := consume.GetHandler()
		if handler == nil {
			return stackErr.Error(fmt.Errorf("consumer handler is nil"))
		}

		if err = handler(ctx, vals); err != nil {
			return stackErr.Error(err)
		}

		return nil
	}
}
