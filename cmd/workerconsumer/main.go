package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/messaging/route"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/otelkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	cfg := config.NewConfig()

	logkit.SetupLogger(cfg)
	validatorkit.SetupValidator(cfg)

	db := provider.NewDatabase(cfg)
	s3Client := provider.NewAWSS3Client(cfg)
	producer := provider.NewKafkaClientProducer(cfg)
	redisClient := provider.NewRedisClient(cfg)
	elasticsearchClient := provider.NewElasticsearchClient(cfg)

	usecases := dependency_injection.SetupUsecases(cfg, db, producer, s3Client, redisClient, elasticsearchClient)

	consumers := dependency_injection.SetupConsumers(cfg, usecases)

	stopTraceProvider := telemetry.InitTraceProvider(cfg)
	defer stopTraceProvider()

	otelkit.ValidateAbleToExportSpan()

	stopLogProvider := telemetry.InitLogProvider(cfg)
	defer stopLogProvider()

	runConsumers(cfg, producer, consumers)
}

func runConsumers(cfg *config.Config, producer *kgo.Client, consumers *dependency_injection.Consumers) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	logkit.Logger.Info("starting worker service")

	route.Setup(ctx, cfg, producer, consumers, wg)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	s := <-terminateSignals
	logkit.Logger.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)

	logkit.Logger.Info("canceling")
	cancel()
	logkit.Logger.Info("canceled")

	logkit.Logger.Info("wait for all consumer to finish processing")
	wg.Wait()
	logkit.Logger.Info("done waiting")

	logkit.Logger.Info("end process of worker")
}
