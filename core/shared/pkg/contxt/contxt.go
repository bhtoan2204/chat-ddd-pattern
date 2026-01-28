package contxt

import (
	"context"
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
