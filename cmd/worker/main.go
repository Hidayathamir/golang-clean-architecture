package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	s3Client := provider.NewAWSS3Client(cfg)
	producer := provider.NewKafkaClientProducer(cfg)
	redisClient := provider.NewRedisClient(cfg)

	usecases := dependency_injection.SetupUsecases(cfg, db, producer, s3Client, redisClient)

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
	wg := &sync.WaitGroup{}

	x.Logger.Info("starting worker service")

	route.Setup(ctx, cfg, consumers, wg)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	s := <-terminateSignals
	x.Logger.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)

	x.Logger.Info("canceling")
	cancel()
	x.Logger.Info("canceled")

	x.Logger.Info("wait for all consumer to finish processing")
	wg.Wait()
	x.Logger.Info("done waiting")

	x.Logger.Info("end process of worker")
}
