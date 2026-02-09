package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/route"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
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
	s3Client := provider.NewS3Client(cfg)
	producer := provider.NewKafkaProducer(cfg)

	usecases := dependency_injection.SetupUsecases(cfg, db, producer, s3Client)

	controllers := dependency_injection.SetupControllers(cfg, usecases)
	middlewares := dependency_injection.SetupMiddlewares(usecases)

	stopTraceProvider, err := telemetry.InitTraceProvider(cfg)
	panicIfErr(err)
	defer stopTraceProvider()

	validateAbleToExportSpan()

	stopLogProvider, err := telemetry.InitLogProvider(cfg)
	panicIfErr(err)
	defer stopLogProvider()

	app := provider.NewFiber(cfg)
	route.Setup(app, controllers, middlewares)

	runHTTPServer(cfg, app)
}

func runHTTPServer(cfg *config.Config, app *fiber.App) {
	// Watch for termination signals so the server can exit gracefully.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()          // Block until a termination signal arrives, then trigger a graceful shutdown.
		err := app.Shutdown() // Stop accepting new requests and wait for in-flight handlers to finish.
		panicIfErr(err)
	}()

	webPort := cfg.GetString(configkey.WebPort)
	fmt.Printf("Go to swagger http://localhost:%s/swagger\n", webPort)
	err := app.Listen(":" + webPort) // Start the HTTP server and block until app shutdown.
	panicIfErr(err)
}

func validateAbleToExportSpan() {
	tracer := otel.Tracer("manual-validation-web")
	_, span := tracer.Start(context.Background(), "startup-check-web")
	span.SetAttributes(attribute.String("check", "success"))
	span.End()

	flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if tp, ok := otel.GetTracerProvider().(*trace.TracerProvider); ok {
		err := tp.ForceFlush(flushCtx)
		if err != nil {
			err = errkit.SetMessage(err, "error export span, wait a little longer, or check is the collector ready")
			x.Logger.WithError(err).Panic()
		}
		x.Logger.Info("Successfully sent manual trace for web")
	}
}

func panicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
