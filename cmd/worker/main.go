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
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func main() {
	viperConfig := config.NewViper()
	x.SetupAll(viperConfig)
	db := config.NewDatabase(viperConfig)
	s3Client := config.NewS3Client(viperConfig)
	producer := config.NewKafkaProducer(viperConfig)

	usecases := config.SetupUsecases(viperConfig, db, producer, s3Client)

	stopTraceProvider, err := telemetry.InitTraceProvider(viperConfig)
	panicIfErr(err)
	defer stopTraceProvider()

	stopLogProvider, err := telemetry.InitLogProvider(viperConfig)
	panicIfErr(err)
	defer stopLogProvider()

	ctx, cancel := context.WithCancel(context.Background())

	x.Logger.Info("Starting worker service")
	go route.SetupUserFollowedConsumer(ctx, viperConfig, usecases)
	go route.SetupImageUploadedConsumer(ctx, viperConfig, usecases)
	go route.SetupImageLikedConsumer(ctx, viperConfig, usecases)
	go route.SetupImageCommentedConsumer(ctx, viperConfig, usecases)
	go route.SetupNotifConsumer(ctx, viperConfig, usecases)

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

func panicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
