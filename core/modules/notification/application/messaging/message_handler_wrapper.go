package messaging

import (
	"context"
	"fmt"
	"go-socket/core/shared/infra/messaging"
	"go-socket/core/shared/pkg/contxt"
	"go-socket/core/shared/pkg/logging"
	"go-socket/core/shared/pkg/stackErr"
)

func (h *messageHandler) processMessage(consume messaging.Consumer) messaging.CallBack {
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

		if err = consume.GetHandler()(ctx, vals); err != nil {
			return stackErr.Error(err)
		}

		return nil
	}
}
