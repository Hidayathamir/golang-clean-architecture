package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/route"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/otelkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

// General API Info
// for swag (https://github.com/swaggo/swag)

//	@title	Golang Clean Architecture

//	@securityDefinitions.apikey	SimpleApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Simple token authorization

func main() {
	cfg := config.NewConfig()

	x.SetupAll(cfg)

	db := provider.NewDatabase(cfg)
	awsS3Client := provider.NewAWSS3Client(cfg)
	producer := provider.NewKafkaClientProducer(cfg)
	redisClient := provider.NewRedisClient(cfg)
	elasticsearchClient := provider.NewelasticsearchClient(cfg)

	usecases := dependency_injection.SetupUsecases(cfg, db, producer, awsS3Client, redisClient, elasticsearchClient)

	controllers := dependency_injection.SetupControllers(cfg, usecases)
	middlewares := dependency_injection.SetupMiddlewares(usecases)

	stopTraceProvider := telemetry.InitTraceProvider(cfg)
	defer stopTraceProvider()

	otelkit.ValidateAbleToExportSpan()

	stopLogProvider := telemetry.InitLogProvider(cfg)
	defer stopLogProvider()

	runHTTPServer(cfg, controllers, middlewares)
}

func runHTTPServer(cfg *config.Config, controllers *dependency_injection.Controllers, middlewares *dependency_injection.Middlewares) {
	app := provider.NewFiber(cfg)
	route.Setup(app, controllers, middlewares)

	// Watch for termination signals so the server can exit gracefully.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()          // Block until a termination signal arrives, then trigger a graceful shutdown.
		err := app.Shutdown() // Stop accepting new requests and wait for in-flight handlers to finish.
		x.PanicIfErr(err)
	}()

	webPort := cfg.GetWebPort()
	fmt.Printf("Go to swagger http://localhost:%s/swagger\n", webPort)
	err := app.Listen(":" + webPort) // Start the HTTP server and block until app shutdown.
	x.PanicIfErr(err)
}
