package route

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/consttopic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/spf13/viper"
)

func ConsumeUserFollowedEventForNotification(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consumer.ConsumeUserFollowedEventForNotification")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewUserConsumer(usecases.UserUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.UserFollowed, consumer.ConsumeUserFollowedEventForNotification)
}

func ConsumeUserFollowedEventForUpdateCount(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consumer.ConsumeUserFollowedEventForUpdateCount")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewUserConsumer(usecases.UserUsecase)
	messaging.ConsumeTopicBatch(ctx, consumerGroup, consttopic.UserFollowed, consumer.ConsumeUserFollowedEventForUpdateCount, 10, 10*time.Second)
}

func SetupImageUploadedConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consumer.ConsumeImageUploadedEvent")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.ImageUploaded, consumer.ConsumeImageUploadedEvent)
}

func SetupImageLikedConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consumer.ConsumeImageLikedEvent")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.ImageLiked, consumer.ConsumeImageLikedEvent)
}

func SetupImageCommentedConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consumer.ConsumeImageCommentedEvent")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.ImageCommented, consumer.ConsumeImageCommentedEvent)
}

func SetupNotifConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consumer.ConsumeNotifEvent")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewNotifConsumer(usecases.NotifUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.Notif, consumer.ConsumeNotifEvent)
}
