package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/otelkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func main() {
	cfg := config.NewConfig()

	logkit.SetupLogger(cfg)
	validatorkit.SetupValidator(cfg)

	db := provider.NewDatabase(cfg)
	producer := provider.NewKafkaClientProducer(cfg)

	// setup outbox repository
	var outboxRepository repository.OutboxRepository
	outboxRepository = repository.NewOutboxRepository(cfg)
	outboxRepository = repository.NewOutboxRepositoryMwLogger(outboxRepository)

	// setup outbox publisher
	var publisher messaging.OutboxPublisher
	publisher = messaging.NewOutboxPublisher(cfg, db, producer, outboxRepository)
	publisher = messaging.NewOutboxPublisherMwLogger(publisher)

	stopTraceProvider := telemetry.InitTraceProvider(cfg)
	defer stopTraceProvider()

	otelkit.ValidateAbleToExportSpan()

	stopLogProvider := telemetry.InitLogProvider(cfg)
	defer stopLogProvider()

	runPublisherLoop(cfg, publisher)
}

func runPublisherLoop(cfg *config.Config, publisher messaging.OutboxPublisher) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	logkit.Logger.Info("starting outbox publisher worker")

	interval := time.Duration(cfg.GetOutboxPollIntervalSeconds()) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				err := publisher.PublishPending(ctx)
				if err != nil {
					logkit.Logger.WithContext(ctx).WithError(err).Error("outbox publish failed")
				}
			case <-ctx.Done():
				// Process remaining pending records one last time before exit
				err := publisher.PublishPending(context.Background())
				if err != nil {
					logkit.Logger.WithError(err).Error("outbox final publish failed")
				}
				return
			}
		}
	}()

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	s := <-terminateSignals
	logkit.Logger.Info("Got one of stop signals, shutting down outbox publisher, SIGNAL NAME :", s)

	logkit.Logger.Info("canceling")
	cancel()
	logkit.Logger.Info("canceled")

	logkit.Logger.Info("wait for all publisher cycle to finish")
	wg.Wait()
	logkit.Logger.Info("done waiting")

	logkit.Logger.Info("end process of outbox publisher")
}
