package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/http/route"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/otelkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
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

	logkit.SetupLogger(cfg)
	validatorkit.SetupValidator(cfg)

	db := provider.NewDatabase(cfg)
	awsS3Client := provider.NewAWSS3Client(cfg)
	redisClient := provider.NewRedisClient(cfg)
	elasticsearchClient := provider.NewElasticsearchClient(cfg)

	repositories := dependency_injection.SetupRepositories(cfg)
	messagings := dependency_injection.SetupMessagings(cfg, repositories)
	caches := dependency_injection.SetupCaches(redisClient)
	searches := dependency_injection.SetupSearches(elasticsearchClient)
	storages := dependency_injection.SetupStorages(cfg, awsS3Client)
	emails := dependency_injection.SetupEmails(cfg)
	hubs := dependency_injection.SetupHubs()

	usecases := dependency_injection.SetupUsecases(cfg, db, caches, messagings, repositories, searches, storages, emails)

	controllers := dependency_injection.SetupControllers(cfg, redisClient, hubs, usecases)
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

	// Forward notifications published by the worker processes to the local sockets.
	go controllers.NotificationSubscriber.Run(ctx)

	go func() {
		<-ctx.Done()          // Block until a termination signal arrives, then trigger a graceful shutdown.
		err := app.Shutdown() // Stop accepting new requests and wait for in-flight handlers to finish.
		errkit.PanicIfErr(err)
	}()

	webPort := cfg.GetWebPort()
	fmt.Printf("Go to swagger http://localhost:%s/swagger\n", webPort)
	err := app.Listen(":" + webPort) // Start the HTTP server and block until app shutdown.
	errkit.PanicIfErr(err)
}
