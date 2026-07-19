package telemetry

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel"
)

type mapCarrier map[string]string

func (m mapCarrier) Get(key string) string {
	return m[key]
}

func (m mapCarrier) Set(key string, value string) {
	m[key] = value
}

func (m mapCarrier) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func InjectTraceContext(ctx context.Context) string {
	carrier := make(mapCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	if len(carrier) == 0 {
		return ""
	}
	b, _ := json.Marshal(carrier)
	return string(b)
}

func ExtractTraceContext(ctx context.Context, traceContext string) context.Context {
	if traceContext == "" {
		return ctx
	}
	carrier := make(mapCarrier)
	json.Unmarshal([]byte(traceContext), &carrier)
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}
