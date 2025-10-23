package route

import (
	_ "github.com/Hidayathamir/golang-clean-architecture/api" // need import for swagger
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Setup(app *fiber.App, controllers *config.Controllers, middlewares *config.Middlewares) {
	setupHomeRoute(app)
	setupSwaggerRoute(app)
	setupGuestRoute(app, controllers, middlewares)
	setupAuthRoute(app, controllers, middlewares)
}

func setupHomeRoute(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("hi")
	})
}

func setupSwaggerRoute(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}

func setupGuestRoute(app *fiber.App, controllers *config.Controllers, middlewares *config.Middlewares) {
	router := app.Group("", middlewares.TraceIDMiddleware)
	router.Post("/api/users", controllers.UserController.Register)
	router.Post("/api/users/_login", controllers.UserController.Login)
}

func setupAuthRoute(app *fiber.App, controllers *config.Controllers, middlewares *config.Middlewares) {
	router := app.Group("", middlewares.TraceIDMiddleware, middlewares.AuthMiddleware)
	router.Delete("/api/users", controllers.UserController.Logout)
	router.Patch("/api/users/_current", controllers.UserController.Update)
	router.Get("/api/users/_current", controllers.UserController.Current)

	router.Get("/api/contacts", controllers.ContactController.List)
	router.Post("/api/contacts", controllers.ContactController.Create)
	router.Put("/api/contacts/:contactId", controllers.ContactController.Update)
	router.Get("/api/contacts/:contactId", controllers.ContactController.Get)
	router.Delete("/api/contacts/:contactId", controllers.ContactController.Delete)

	router.Get("/api/contacts/:contactId/addresses", controllers.AddressController.List)
	router.Post("/api/contacts/:contactId/addresses", controllers.AddressController.Create)
	router.Put("/api/contacts/:contactId/addresses/:addressId", controllers.AddressController.Update)
	router.Get("/api/contacts/:contactId/addresses/:addressId", controllers.AddressController.Get)
	router.Delete("/api/contacts/:contactId/addresses/:addressId", controllers.AddressController.Delete)

	router.Post("/api/todos", controllers.TodoController.Create)
	router.Get("/api/todos", controllers.TodoController.List)
	router.Get("/api/todos/:todoId", controllers.TodoController.Get)
	router.Put("/api/todos/:todoId", controllers.TodoController.Update)
	router.Delete("/api/todos/:todoId", controllers.TodoController.Delete)
	router.Patch("/api/todos/:todoId/_complete", controllers.TodoController.Complete)
}
