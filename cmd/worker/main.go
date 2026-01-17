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
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/consttopic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/spf13/viper"
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
	go RunUserFollowedConsumer(ctx, viperConfig, usecases)
	go RunImageUploadedConsumer(ctx, viperConfig, usecases)
	go RunImageLikedConsumer(ctx, viperConfig, usecases)
	go RunImageCommentedConsumer(ctx, viperConfig, usecases)
	go RunNotifConsumer(ctx, viperConfig, usecases)

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

func RunAddressConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup address consumer")
	addressConsumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	addressHandler := messaging.NewAddressConsumer(usecases.AddressUsecase)
	messaging.ConsumeTopic(ctx, addressConsumerGroup, "addresses", addressHandler.Consume)
}

func RunContactConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup contact consumer")
	contactConsumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	contactHandler := messaging.NewContactConsumer(usecases.ContactUsecase)
	messaging.ConsumeTopic(ctx, contactConsumerGroup, "contacts", contactHandler.Consume)
}

func RunUserFollowedConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consttopic.UserFollowed consumer")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewUserConsumer(usecases.UserUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.UserFollowed, consumer.ConsumeUserFollowedEvent)
}

func RunTodoCommandConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup todo command consumer")
	todoCommandGroup := config.NewKafkaConsumerGroup(viperConfig)
	commandHandler := messaging.NewTodoCommandConsumer(usecases.TodoUsecase)
	messaging.ConsumeTopic(ctx, todoCommandGroup, "todo-commands", commandHandler.Consume)
}

func RunTodoCompletionConsumer(ctx context.Context, viperConfig *viper.Viper) {
	x.Logger.Info("setup todo completion consumer")
	todoCompletionGroup := config.NewKafkaConsumerGroup(viperConfig)
	completionHandler := messaging.NewTodoCompletionConsumer()
	messaging.ConsumeTopic(ctx, todoCompletionGroup, "todos", completionHandler.Consume)
}

func RunImageUploadedConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consttopic.ImageUploaded consumer")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.ImageUploaded, consumer.ConsumeImageUploadedEvent)
}

func RunImageLikedConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consttopic.ImageLiked consumer")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.ImageLiked, consumer.ConsumeImageLikedEvent)
}

func RunImageCommentedConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consttopic.ImageCommented consumer")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.ImageCommented, consumer.ConsumeImageCommentedEvent)
}

func RunNotifConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consttopic.ImageCommented consumer")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewNotifConsumer(usecases.NotifUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.Notif, consumer.ConsumeNotifEvent)
}

func panicIfErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
