package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/messaging/route"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	cfg := config.NewConfig()

	x.SetupAll(cfg)

	db := provider.NewDatabase(cfg)
	s3Client := provider.NewS3Client(cfg)
	producer := provider.NewKafkaProducer(cfg)

	usecases := dependency_injection.SetupUsecases(cfg, db, producer, s3Client)

	consumers := dependency_injection.SetupConsumers(cfg, usecases)

	stopTraceProvider, err := telemetry.InitTraceProvider(cfg)
	panicIfErr(err)
	defer stopTraceProvider()

	validateAbleToExportSpan()

	stopLogProvider, err := telemetry.InitLogProvider(cfg)
	panicIfErr(err)
	defer stopLogProvider()

	runConsumers(cfg, consumers)
}

func runConsumers(cfg *config.Config, consumers *dependency_injection.Consumers) {
	ctx, cancel := context.WithCancel(context.Background())

	x.Logger.Info("Starting worker service")
	go route.ConsumeUserFollowedEventForNotification(ctx, cfg, consumers)
	go route.ConsumeUserFollowedEventForUpdateCount(ctx, cfg, consumers)
	go route.SetupImageUploadedConsumer(ctx, cfg, consumers)
	go route.ConsumeImageLikedEventForNotification(ctx, cfg, consumers)
	go route.ConsumeImageLikedEventForUpdateCount(ctx, cfg, consumers)
	go route.ConsumeImageCommentedEventForNotification(ctx, cfg, consumers)
	go route.ConsumeImageCommentedEventForUpdateCount(ctx, cfg, consumers)
	go route.SetupNotifConsumer(ctx, cfg, consumers)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	s := <-terminateSignals
	x.Logger.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)

	x.Logger.Info("canceling")
	cancel()
	x.Logger.Info("canceled")

	x.Logger.Info("wait for all consumer to finish processing, 5 second")
	time.Sleep(5 * time.Second)
	x.Logger.Info("done waiting")

	x.Logger.Info("end process of worker")
}

func validateAbleToExportSpan() {
	tracer := otel.Tracer("manual-validation-worker")
	_, span := tracer.Start(context.Background(), "startup-check-worker")
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
		x.Logger.Info("Successfully sent manual trace for worker")
	}
}

func panicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
