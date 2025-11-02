package middleware

import (
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/gofiber/fiber/v2"
)

func NewTraceID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		traceID := telemetry.GetTraceID(c.UserContext())
		if traceID == "" {
			return c.Next()
		}

		c.Set("X-Trace-ID", traceID)

		return c.Next()
	}
}
