package middleware

import (
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/traceidctx"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func NewTraceID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := c.Get("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		ctx := traceidctx.Set(c.UserContext(), traceID)
		c.SetUserContext(ctx)

		c.Set("X-Trace-ID", traceID)

		return c.Next()
	}
}
