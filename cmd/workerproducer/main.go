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

	// setup outbox producer
	var outboxProducer messaging.OutboxProducer
	outboxProducer = messaging.NewOutboxProducer(cfg, db, producer, outboxRepository)
	outboxProducer = messaging.NewOutboxProducerMwLogger(outboxProducer)

	stopTraceProvider := telemetry.InitTraceProvider(cfg)
	defer stopTraceProvider()

	otelkit.ValidateAbleToExportSpan()

	stopLogProvider := telemetry.InitLogProvider(cfg)
	defer stopLogProvider()

	runProducerLoop(cfg, outboxProducer)
}

func runProducerLoop(cfg *config.Config, producer messaging.OutboxProducer) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	logkit.Logger.Info("starting outbox producer worker")

	interval := time.Duration(cfg.GetOutboxPollIntervalSeconds()) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				err := producer.ProducePending(ctx)
				if err != nil {
					logkit.Logger.WithContext(ctx).WithError(err).Error("outbox produce failed")
				}
			case <-ctx.Done():
				// Process remaining pending records one last time before exit
				err := producer.ProducePending(context.Background())
				if err != nil {
					logkit.Logger.WithError(err).Error("outbox final produce failed")
				}
				return
			}
		}
	}()

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	s := <-terminateSignals
	logkit.Logger.Info("Got one of stop signals, shutting down outbox producer, SIGNAL NAME :", s)

	logkit.Logger.Info("canceling")
	cancel()
	logkit.Logger.Info("canceled")

	logkit.Logger.Info("wait for all producer cycle to finish")
	wg.Wait()
	logkit.Logger.Info("done waiting")

	logkit.Logger.Info("end process of outbox producer")
}
