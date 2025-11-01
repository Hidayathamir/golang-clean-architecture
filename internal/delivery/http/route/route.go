package route

import (
	_ "github.com/Hidayathamir/golang-clean-architecture/api" // need import for swagger
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Setup(app *fiber.App, controllers *config.Controllers, middlewares *config.Middlewares) {
	root := app.Group("")

	setupHomeRoute(root)
	setupSwaggerRoute(root)

	api := root.Group("/api", middlewares.OtelFiberMiddleware, middlewares.TraceIDMiddleware)

	setupGuestRoute(api, controllers)

	authenticated := api.Group("", middlewares.AuthMiddleware)
	setupAuthRoute(authenticated, controllers)
}

func setupHomeRoute(router fiber.Router) {
	router.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("hi")
	})
}

func setupSwaggerRoute(router fiber.Router) {
	router.Get("/swagger/*", swagger.HandlerDefault)
}

func setupGuestRoute(router fiber.Router, controllers *config.Controllers) {
	users := router.Group("/users")
	{
		users.Post("", controllers.UserController.Register)
		users.Post("/_login", controllers.UserController.Login)
	}
}

func setupAuthRoute(router fiber.Router, controllers *config.Controllers) {
	users := router.Group("/users")
	{
		users.Delete("", controllers.UserController.Logout)
		users.Patch("/_current", controllers.UserController.Update)
		users.Get("/_current", controllers.UserController.Current)
	}

	contacts := router.Group("/contacts")
	{
		contacts.Get("", controllers.ContactController.List)
		contacts.Post("", controllers.ContactController.Create)
		contacts.Put("/:contactId", controllers.ContactController.Update)
		contacts.Get("/:contactId", controllers.ContactController.Get)
		contacts.Delete("/:contactId", controllers.ContactController.Delete)
	}

	addresses := contacts.Group("/:contactId/addresses")
	{
		addresses.Get("", controllers.AddressController.List)
		addresses.Post("", controllers.AddressController.Create)
		addresses.Put("/:addressId", controllers.AddressController.Update)
		addresses.Get("/:addressId", controllers.AddressController.Get)
		addresses.Delete("/:addressId", controllers.AddressController.Delete)
	}

	todos := router.Group("/todos")
	{
		todos.Post("", controllers.TodoController.Create)
		todos.Get("", controllers.TodoController.List)
		todos.Get("/:todoId", controllers.TodoController.Get)
		todos.Put("/:todoId", controllers.TodoController.Update)
		todos.Delete("/:todoId", controllers.TodoController.Delete)
		todos.Patch("/:todoId/_complete", controllers.TodoController.Complete)
	}
}
