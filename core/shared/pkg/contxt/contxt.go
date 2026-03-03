package contxt

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey string

const ctxKeyRequestID = ctxKey("request_id")

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKeyRequestID, id)
}

func RequestIDFromCtx(ctx context.Context) string {
	v := ctx.Value(ctxKeyRequestID)
	if v == nil {
		return ""
	}

	if val, ok := v.(string); ok {
		return val
	}

	return ""
}

func SetRequestID(ctx context.Context) context.Context {
	reqID := RequestIDFromCtx(ctx)
	if reqID == "" {
		reqID = uuid.NewString()
	}
	return WithRequestID(ctx, reqID)
}
