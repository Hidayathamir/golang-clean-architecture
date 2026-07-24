package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/idempotencyusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func main() {
	cfg := config.NewConfig()

	logkit.SetupLogger(cfg)
	validatorkit.SetupValidator(cfg)

	db := provider.NewDatabase(cfg)

	var idempotencyRepo repository.IdempotencyRepository
	idempotencyRepo = repository.NewIdempotencyRepository(cfg)
	idempotencyRepo = repository.NewIdempotencyRepositoryMwLogger(idempotencyRepo)

	var idempotencyUsecase idempotencyusecase.IdempotencyUsecase
	idempotencyUsecase = idempotencyusecase.NewIdempotencyUsecase(cfg, db, idempotencyRepo)
	idempotencyUsecase = idempotencyusecase.NewIdempotencyUsecaseMwLogger(idempotencyUsecase)

	stopTraceProvider := telemetry.InitTraceProvider(cfg)
	defer stopTraceProvider()

	stopLogProvider := telemetry.InitLogProvider(cfg)
	defer stopLogProvider()

	runCleaner(cfg, idempotencyUsecase)
}

func runCleaner(cfg *config.Config, usecase idempotencyusecase.IdempotencyUsecase) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	logkit.Logger.Info("starting idempotency cleanup worker")

	interval := time.Duration(cfg.GetIdempotencyCleanupIntervalSeconds()) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				deleted, err := usecase.DeleteOlderThan(ctx, 7*24*time.Hour)
				if err != nil {
					logkit.Logger.WithContext(ctx).WithError(err).Error("idempotency cleanup failed")
				} else if deleted > 0 {
					logkit.Logger.WithContext(ctx).WithField("deleted", deleted).Info("idempotency records cleaned")
				}
			case <-ctx.Done():
				deleted, err := usecase.DeleteOlderThan(context.Background(), 7*24*time.Hour)
				if err != nil {
					logkit.Logger.WithError(err).Error("idempotency final cleanup failed")
				} else if deleted > 0 {
					logkit.Logger.WithField("deleted", deleted).Info("idempotency final cleanup done")
				}
				return
			}
		}
	}()

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	s := <-terminateSignals
	logkit.Logger.Info("Got one of stop signals, shutting down idempotency cleaner, SIGNAL NAME :", s)

	logkit.Logger.Info("canceling")
	cancel()
	logkit.Logger.Info("canceled")

	logkit.Logger.Info("wait for all cleanup cycle to finish")
	wg.Wait()
	logkit.Logger.Info("done waiting")

	logkit.Logger.Info("end process of idempotency cleaner")
}
