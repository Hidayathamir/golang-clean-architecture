package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/route"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/l"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
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
	l.SetupLogger(viperConfig)
	db := config.NewDatabase(viperConfig)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	producer := config.NewKafkaProducer(viperConfig)

	usecases := config.SetupUsecases(viperConfig, db, app, validate, producer)
	controllers := config.SetupControllers(viperConfig, usecases)
	middlewares := config.SetupMiddlewares(usecases)

	stopTraceProvider, err := telemetry.InitTraceProvider(viperConfig)
	panicIfErr(err)
	defer stopTraceProvider()

	stopLogProvider, err := telemetry.InitLogProvider(viperConfig)
	panicIfErr(err)
	defer stopLogProvider()

	route.Setup(app, controllers, middlewares)

	runHTTPServer(viperConfig, app)
}

func runHTTPServer(viperConfig *viper.Viper, app *fiber.App) {
	// Watch for termination signals so the server can exit gracefully.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()          // Block until a termination signal arrives, then trigger a graceful shutdown.
		err := app.Shutdown() // Stop accepting new requests and wait for in-flight handlers to finish.
		panicIfErr(err)
	}()

	webPort := viperConfig.GetString(configkey.WebPort)
	fmt.Printf("Go to swagger http://localhost:%s/swagger\n", webPort)
	err := app.Listen(":" + webPort) // Start the HTTP server and block until app shutdown.
	panicIfErr(err)
}

func panicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
