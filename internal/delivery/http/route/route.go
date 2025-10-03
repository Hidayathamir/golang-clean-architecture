package route

import (
	_ "github.com/Hidayathamir/golang-clean-architecture/api" // need import for swagger
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Setup(
	app *fiber.App,
	userController *http.UserController,
	contactController *http.ContactController,
	addressController *http.AddressController,
	authMiddleware fiber.Handler,
	traceIDMiddleware fiber.Handler,
) {
	setupHomeRoute(app)
	setupSwaggerRoute(app)
	setupGuestRoute(
		app,
		userController,
		traceIDMiddleware,
	)
	setupAuthRoute(
		app,
		userController,
		contactController,
		addressController,
		authMiddleware,
		traceIDMiddleware,
	)
}

func setupHomeRoute(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("hi")
	})
}

func setupSwaggerRoute(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}

func setupGuestRoute(
	app *fiber.App,
	userController *http.UserController,
	traceIDMiddleware fiber.Handler,
) {
	router := app.Group("", traceIDMiddleware)
	router.Post("/api/users", userController.Register)
	router.Post("/api/users/_login", userController.Login)
}

func setupAuthRoute(
	app *fiber.App,
	userController *http.UserController,
	contactController *http.ContactController,
	addressController *http.AddressController,
	authMiddleware fiber.Handler,
	traceIDMiddleware fiber.Handler,
) {
	router := app.Group("", traceIDMiddleware, authMiddleware)
	router.Delete("/api/users", userController.Logout)
	router.Patch("/api/users/_current", userController.Update)
	router.Get("/api/users/_current", userController.Current)

	router.Get("/api/contacts", contactController.List)
	router.Post("/api/contacts", contactController.Create)
	router.Put("/api/contacts/:contactId", contactController.Update)
	router.Get("/api/contacts/:contactId", contactController.Get)
	router.Delete("/api/contacts/:contactId", contactController.Delete)

	router.Get("/api/contacts/:contactId/addresses", addressController.List)
	router.Post("/api/contacts/:contactId/addresses", addressController.Create)
	router.Put("/api/contacts/:contactId/addresses/:addressId", addressController.Update)
	router.Get("/api/contacts/:contactId/addresses/:addressId", addressController.Get)
	router.Delete("/api/contacts/:contactId/addresses/:addressId", addressController.Delete)
}
