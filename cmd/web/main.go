package main

import (
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/route"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
)

// General API Info
// for swag (https://github.com/swaggo/swag)

//	@title	Golang Clean Architecture

//	@securityDefinitions.apikey	SimpleApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Simple token authorization

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	producer := config.NewKafkaProducer(viperConfig, log)

	usecases := config.SetupUsecases(db, app, log, validate, viperConfig, producer)

	userController := http.NewUserController(usecases.UserUsecase, log)
	contactController := http.NewContactController(usecases.ContactUsecase, log)
	addressController := http.NewAddressController(usecases.AddressUsecase, log)

	authMiddleware := middleware.NewAuth(usecases.UserUsecase)
	traceIDMiddleware := middleware.NewTraceID()

	route.Setup(
		app,
		userController,
		contactController,
		addressController,
		authMiddleware,
		traceIDMiddleware,
	)

	webPort := viperConfig.GetString(configkey.WebPort)
	fmt.Printf("Go to swagger http://localhost:%s/swagger\n", webPort)
	err := app.Listen(":" + webPort)
	if err != nil {
		log.Panicf("Failed to start server: %v", err)
	}
}
