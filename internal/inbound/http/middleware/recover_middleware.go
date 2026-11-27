package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// NewRecover converts a panic in a downstream handler into an error so it flows
// through the central fiber ErrorHandler instead of dropping the connection.
func NewRecover() fiber.Handler {
	return recover.New(recover.Config{EnableStackTrace: true})
}
