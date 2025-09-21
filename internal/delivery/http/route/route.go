package route

import (
	_ "golang-clean-architecture/api" // need import for swagger
	"golang-clean-architecture/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type RouteConfig struct {
	App               *fiber.App
	UserController    *http.UserController
	ContactController *http.ContactController
	AddressController *http.AddressController
	AuthMiddleware    fiber.Handler
	TraceIDMiddleware fiber.Handler
}

func (c *RouteConfig) Setup() {
	c.SetupHomeRoute()
	c.SetupSwaggerRoute()
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupHomeRoute() {
	c.App.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("hi")
	})
}

func (c *RouteConfig) SetupSwaggerRoute() {
	c.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (c *RouteConfig) SetupGuestRoute() {
	c.App.Use(c.TraceIDMiddleware)

	c.App.Post("/api/users", c.UserController.Register)
	c.App.Post("/api/users/_login", c.UserController.Login)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.TraceIDMiddleware)
	c.App.Use(c.AuthMiddleware)

	c.App.Delete("/api/users", c.UserController.Logout)
	c.App.Patch("/api/users/_current", c.UserController.Update)
	c.App.Get("/api/users/_current", c.UserController.Current)

	c.App.Get("/api/contacts", c.ContactController.List)
	c.App.Post("/api/contacts", c.ContactController.Create)
	c.App.Put("/api/contacts/:contactId", c.ContactController.Update)
	c.App.Get("/api/contacts/:contactId", c.ContactController.Get)
	c.App.Delete("/api/contacts/:contactId", c.ContactController.Delete)

	c.App.Get("/api/contacts/:contactId/addresses", c.AddressController.List)
	c.App.Post("/api/contacts/:contactId/addresses", c.AddressController.Create)
	c.App.Put("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Update)
	c.App.Get("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Get)
	c.App.Delete("/api/contacts/:contactId/addresses/:addressId", c.AddressController.Delete)
}
