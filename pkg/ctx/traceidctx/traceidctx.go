package traceidctx

import (
	"context"
)

type ctxKey string

const traceIDKey ctxKey = "trace_id"

func Set(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

func Get(ctx context.Context) string {
	if v := ctx.Value(traceIDKey); v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
