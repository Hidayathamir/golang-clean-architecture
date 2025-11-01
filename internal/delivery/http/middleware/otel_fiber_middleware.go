package middleware

import (
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
)

func NewOtelFiberMiddleware() fiber.Handler {
	return otelfiber.Middleware(
		otelfiber.WithSpanNameFormatter(
			func(ctx *fiber.Ctx) string {
				path := ctx.Path()
				if route := ctx.Route(); route != nil && route.Path != "" {
					path = route.Path
				}

				if path == "" {
					path = "/"
				}

				return ctx.Method() + " " + path
			},
		),
	)
}
