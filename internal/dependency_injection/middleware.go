package dependency_injection

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/gofiber/fiber/v2"
)

type Middlewares struct {
	AuthMiddleware      fiber.Handler
	TraceIDMiddleware   fiber.Handler
	OtelFiberMiddleware fiber.Handler
}

func SetupMiddlewares(usecases *Usecases) *Middlewares {
	authMiddleware := middleware.NewAuth(usecases.UserUsecase)
	traceIDMiddleware := middleware.NewTraceID()
	otelFiberMiddleware := middleware.NewOtelFiberMiddleware()

	return &Middlewares{
		AuthMiddleware:      authMiddleware,
		TraceIDMiddleware:   traceIDMiddleware,
		OtelFiberMiddleware: otelFiberMiddleware,
	}
}
