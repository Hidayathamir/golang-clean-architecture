package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	sc := span.SpanContext()
	if !sc.IsValid() || !sc.TraceID().IsValid() {
		return ""
	}

	return sc.TraceID().String()
}

func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return ""
	}

	sc := span.SpanContext()
	if !sc.IsValid() || !sc.SpanID().IsValid() {
		return ""
	}

	return sc.SpanID().String()
}
