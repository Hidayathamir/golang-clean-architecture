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

	usecases := config.SetupUsecases(viperConfig, log, db, app, validate, producer)

	ctx, cancel := context.WithCancel(context.Background())

	log.Info("Starting worker service")
	go RunUserConsumer(ctx, log, viperConfig, usecases)
	go RunContactConsumer(ctx, log, viperConfig, usecases)
	go RunAddressConsumer(ctx, log, viperConfig, usecases)
	go RunTodoCommandConsumer(ctx, log, viperConfig, usecases)
	go RunTodoCompletionConsumer(ctx, log, viperConfig)

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

func RunTodoCommandConsumer(ctx context.Context, log *logrus.Logger, viperConfig *viper.Viper, usecases *config.Usecases) {
	log.Info("setup todo command consumer")
	todoCommandGroup := config.NewKafkaConsumerGroup(viperConfig, log)
	commandHandler := messaging.NewTodoCommandConsumer(usecases.TodoUsecase, log)
	messaging.ConsumeTopic(ctx, todoCommandGroup, "todo-commands", log, commandHandler.Consume)
}

func RunTodoCompletionConsumer(ctx context.Context, log *logrus.Logger, viperConfig *viper.Viper) {
	log.Info("setup todo completion consumer")
	todoCompletionGroup := config.NewKafkaConsumerGroup(viperConfig, log)
	completionHandler := messaging.NewTodoCompletionConsumer(log)
	messaging.ConsumeTopic(ctx, todoCompletionGroup, "todos", log, completionHandler.Consume)
}
