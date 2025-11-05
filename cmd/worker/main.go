package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/l"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/spf13/viper"
)

func main() {
	viperConfig := config.NewViper()
	l.SetupLogger(viperConfig)
	db := config.NewDatabase(viperConfig)
	validate := config.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)
	producer := config.NewKafkaProducer(viperConfig)

	usecases := config.SetupUsecases(viperConfig, db, app, validate, producer)

	stopTraceProvider, err := telemetry.InitTraceProvider(viperConfig)
	panicIfErr(err)
	defer stopTraceProvider()

	stopLogProvider, err := telemetry.InitLogProvider(viperConfig)
	panicIfErr(err)
	defer stopLogProvider()

	ctx, cancel := context.WithCancel(context.Background())

	l.Logger.Info("Starting worker service")
	go RunUserConsumer(ctx, viperConfig, usecases)
	go RunContactConsumer(ctx, viperConfig, usecases)
	go RunAddressConsumer(ctx, viperConfig, usecases)
	go RunTodoCommandConsumer(ctx, viperConfig, usecases)
	go RunTodoCompletionConsumer(ctx, viperConfig)

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGTERM)

	s := <-terminateSignals
	l.Logger.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
	cancel()

	time.Sleep(5 * time.Second) // wait for all consumers to finish processing
}

func RunAddressConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	l.Logger.Info("setup address consumer")
	addressConsumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	addressHandler := messaging.NewAddressConsumer(usecases.AddressUsecase)
	messaging.ConsumeTopic(ctx, addressConsumerGroup, "addresses", addressHandler.Consume)
}

func RunContactConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	l.Logger.Info("setup contact consumer")
	contactConsumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	contactHandler := messaging.NewContactConsumer(usecases.ContactUsecase)
	messaging.ConsumeTopic(ctx, contactConsumerGroup, "contacts", contactHandler.Consume)
}

func RunUserConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	l.Logger.Info("setup user consumer")
	userConsumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	userHandler := messaging.NewUserConsumer(usecases.UserUsecase)
	messaging.ConsumeTopic(ctx, userConsumerGroup, "users", userHandler.Consume)
}

func RunTodoCommandConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	l.Logger.Info("setup todo command consumer")
	todoCommandGroup := config.NewKafkaConsumerGroup(viperConfig)
	commandHandler := messaging.NewTodoCommandConsumer(usecases.TodoUsecase)
	messaging.ConsumeTopic(ctx, todoCommandGroup, "todo-commands", commandHandler.Consume)
}

func RunTodoCompletionConsumer(ctx context.Context, viperConfig *viper.Viper) {
	l.Logger.Info("setup todo completion consumer")
	todoCompletionGroup := config.NewKafkaConsumerGroup(viperConfig)
	completionHandler := messaging.NewTodoCompletionConsumer()
	messaging.ConsumeTopic(ctx, todoCompletionGroup, "todos", completionHandler.Consume)
}

func panicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
