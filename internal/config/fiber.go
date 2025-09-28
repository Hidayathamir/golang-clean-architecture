package config

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(config *viper.Viper) *fiber.App {
	var app = fiber.New(fiber.Config{
		AppName:      config.GetString(configkey.AppName),
		ErrorHandler: response.Error,
		Prefork:      config.GetBool(configkey.WebPrefork),
	})

	return app
}
