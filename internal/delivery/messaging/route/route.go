package route

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/consttopic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/spf13/viper"
)

func SetupUserFollowedConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consttopic.UserFollowed consumer")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewUserConsumer(usecases.UserUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.UserFollowed, consumer.ConsumeUserFollowedEvent)
}

func SetupImageUploadedConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consttopic.ImageUploaded consumer")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.ImageUploaded, consumer.ConsumeImageUploadedEvent)
}

func SetupImageLikedConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consttopic.ImageLiked consumer")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.ImageLiked, consumer.ConsumeImageLikedEvent)
}

func SetupImageCommentedConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consttopic.ImageCommented consumer")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.ImageCommented, consumer.ConsumeImageCommentedEvent)
}

func SetupNotifConsumer(ctx context.Context, viperConfig *viper.Viper, usecases *config.Usecases) {
	x.Logger.Info("setup consttopic.ImageCommented consumer")
	consumerGroup := config.NewKafkaConsumerGroup(viperConfig)
	consumer := messaging.NewNotifConsumer(usecases.NotifUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, consttopic.Notif, consumer.ConsumeNotifEvent)
}
