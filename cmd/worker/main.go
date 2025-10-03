package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/messaging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	producer := config.NewKafkaProducer(viperConfig, log)

	usecases := config.SetupUsecases(db, app, log, validate, viperConfig, producer)

	ctx, cancel := context.WithCancel(context.Background())

	log.Info("Starting worker service")
	go RunUserConsumer(ctx, log, viperConfig, usecases)
	go RunContactConsumer(ctx, log, viperConfig, usecases)
	go RunAddressConsumer(ctx, log, viperConfig, usecases)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	s := <-terminateSignals
	log.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
	cancel()

	time.Sleep(5 * time.Second) // wait for all consumers to finish processing
}

func RunAddressConsumer(ctx context.Context, log *logrus.Logger, viperConfig *viper.Viper, usecases *config.Usecases) {
	log.Info("setup address consumer")
	addressConsumerGroup := config.NewKafkaConsumerGroup(viperConfig, log)
	addressHandler := messaging.NewAddressConsumer(usecases.AddressUsecase, log)
	messaging.ConsumeTopic(ctx, addressConsumerGroup, "addresses", log, addressHandler.Consume)
}

func RunContactConsumer(ctx context.Context, log *logrus.Logger, viperConfig *viper.Viper, usecases *config.Usecases) {
	log.Info("setup contact consumer")
	contactConsumerGroup := config.NewKafkaConsumerGroup(viperConfig, log)
	contactHandler := messaging.NewContactConsumer(usecases.ContactUsecase, log)
	messaging.ConsumeTopic(ctx, contactConsumerGroup, "contacts", log, contactHandler.Consume)
}

func RunUserConsumer(ctx context.Context, log *logrus.Logger, viperConfig *viper.Viper, usecases *config.Usecases) {
	log.Info("setup user consumer")
	userConsumerGroup := config.NewKafkaConsumerGroup(viperConfig, log)
	userHandler := messaging.NewUserConsumer(usecases.UserUsecase, log)
	messaging.ConsumeTopic(ctx, userConsumerGroup, "users", log, userHandler.Consume)
}
