package provider

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/gofiber/fiber/v2"
)

func NewFiber(cfg *config.Config) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      cfg.GetString(configkey.AppName),
		ErrorHandler: response.Error,
		Prefork:      cfg.GetBool(configkey.WebPrefork),
		BodyLimit:    15 * 1024 * 1024, // 15MB
	})

	return app
}
