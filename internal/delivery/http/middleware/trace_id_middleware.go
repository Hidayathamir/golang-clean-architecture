package middleware

import (
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/traceidctx"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
)

func NewTraceID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		span := trace.SpanFromContext(c.UserContext())
		if span == nil {
			return c.Next()
		}

		sc := span.SpanContext()
		if !sc.IsValid() || !sc.TraceID().IsValid() {
			return c.Next()
		}

		traceID := sc.TraceID().String()

		ctx := traceidctx.Set(c.UserContext(), traceID)
		c.SetUserContext(ctx)

		c.Set("X-Trace-ID", traceID)

		return c.Next()
	}
}
