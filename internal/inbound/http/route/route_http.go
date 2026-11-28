package route

import (
	_ "github.com/Hidayathamir/golang-clean-architecture/api" // need import for swagger
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func Setup(app *fiber.App, controllers *dependency_injection.Controllers, middlewares *dependency_injection.Middlewares) {
	root := app.Group("")

	root.Use(middlewares.RecoverMiddleware)

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

func setupGuestRoute(router fiber.Router, controllers *dependency_injection.Controllers) {
	users := router.Group("/users")
	{
		users.Post("", controllers.UserController.Register)
		users.Post("/_login", controllers.UserController.Login)
	}
}

func setupAuthRoute(router fiber.Router, controllers *dependency_injection.Controllers) {
	users := router.Group("/users")
	{
		users.Patch("/_current", controllers.UserController.Update)
		users.Get("/_current", controllers.UserController.Current)
		users.Post("/_follow", controllers.UserController.Follow)
		users.Get("/:id/images", controllers.ImageController.GetUserImages)
		users.Get("/:id/followers", controllers.UserController.GetFollowers)
		users.Get("/:id/following", controllers.UserController.GetFollowings)
		users.Get("/:id", controllers.UserController.GetProfile)
	}

	images := router.Group("/images")
	{
		images.Post("", controllers.ImageController.Upload)
		images.Post("/_like", controllers.ImageController.Like)
		images.Post("/_comment", controllers.ImageController.Comment)
		images.Get("/:imageId", controllers.ImageController.GetImage)
		images.Delete("/:imageId", controllers.ImageController.DeleteImage)
		images.Get("/:imageId/likes", controllers.ImageController.GetLike)
		images.Get("/:imageId/comments", controllers.ImageController.GetComment)
	}

	router.Get("/feed", controllers.ImageController.GetFeed)
	router.Get("/explore", controllers.ImageController.Explore)
}
