package provider

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/gofiber/fiber/v2"
)

func NewFiber(cfg *config.Config) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      cfg.GetAppName(),
		ErrorHandler: response.Error,
		Prefork:      cfg.GetWebPrefork(),
		BodyLimit:    15 * 1024 * 1024, // 15MB
	})

	return app
}
