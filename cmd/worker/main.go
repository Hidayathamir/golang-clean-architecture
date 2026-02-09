package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/messaging/route"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/otelkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func main() {
	cfg := config.NewConfig()

	x.SetupAll(cfg)

	db := provider.NewDatabase(cfg)
	s3Client := provider.NewS3Client(cfg)
	producer := provider.NewKafkaProducer(cfg)

	usecases := dependency_injection.SetupUsecases(cfg, db, producer, s3Client)

	consumers := dependency_injection.SetupConsumers(cfg, usecases)

	stopTraceProvider := telemetry.InitTraceProvider(cfg)
	defer stopTraceProvider()

	otelkit.ValidateAbleToExportSpan()

	stopLogProvider := telemetry.InitLogProvider(cfg)
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
