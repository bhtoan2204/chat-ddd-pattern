package xtracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return tracer.Start(ctx, name, opts...)
}
