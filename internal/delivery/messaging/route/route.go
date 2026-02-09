package route

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/consumergroup"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func ConsumeUserFollowedEventForNotification(ctx context.Context, cfg *config.Config, usecases *dependency_injection.Usecases) {
	x.Logger.Info("setup consumer.ConsumeUserFollowedEventForNotification")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.UserFollowedGroup1)
	consumer := messaging.NewUserConsumer(usecases.UserUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, topic.UserFollowed, consumer.ConsumeUserFollowedEventForNotification)
}

func ConsumeUserFollowedEventForUpdateCount(ctx context.Context, cfg *config.Config, usecases *dependency_injection.Usecases) {
	x.Logger.Info("setup consumer.ConsumeUserFollowedEventForUpdateCount")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.UserFollowedGroup2)
	consumer := messaging.NewUserConsumer(usecases.UserUsecase)
	messaging.ConsumeTopicBatch(ctx, consumerGroup, topic.UserFollowed, consumer.ConsumeUserFollowedEventForUpdateCount, 10, 10*time.Second)
}

func SetupImageUploadedConsumer(ctx context.Context, cfg *config.Config, usecases *dependency_injection.Usecases) {
	x.Logger.Info("setup consumer.ConsumeImageUploadedEvent")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.ImageUploadedGroup1)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, topic.ImageUploaded, consumer.ConsumeImageUploadedEvent)
}

func ConsumeImageLikedEventForNotification(ctx context.Context, cfg *config.Config, usecases *dependency_injection.Usecases) {
	x.Logger.Info("setup consumer.ConsumeImageLikedEventForNotification")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.ImageLikedGroup1)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, topic.ImageLiked, consumer.ConsumeImageLikedEventForNotification)
}

func ConsumeImageLikedEventForUpdateCount(ctx context.Context, cfg *config.Config, usecases *dependency_injection.Usecases) {
	x.Logger.Info("setup consumer.ConsumeImageLikedEventForUpdateCount")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.ImageLikedGroup2)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopicBatch(ctx, consumerGroup, topic.ImageLiked, consumer.ConsumeImageLikedEventForUpdateCount, 10, 10*time.Second)
}

func ConsumeImageCommentedEventForNotification(ctx context.Context, cfg *config.Config, usecases *dependency_injection.Usecases) {
	x.Logger.Info("setup consumer.ConsumeImageCommentedEventForNotification")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.ImageCommentedGroup1)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, topic.ImageCommented, consumer.ConsumeImageCommentedEventForNotification)
}

func ConsumeImageCommentedEventForUpdateCount(ctx context.Context, cfg *config.Config, usecases *dependency_injection.Usecases) {
	x.Logger.Info("setup consumer.ConsumeImageCommentedEventForUpdateCount")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.ImageCommentedGroup2)
	consumer := messaging.NewImageConsumer(usecases.ImageUsecase)
	messaging.ConsumeTopicBatch(ctx, consumerGroup, topic.ImageCommented, consumer.ConsumeImageCommentedEventForUpdateCount, 10, 18*time.Second)
}

func SetupNotifConsumer(ctx context.Context, cfg *config.Config, usecases *dependency_injection.Usecases) {
	x.Logger.Info("setup consumer.ConsumeNotifEvent")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.NotifGroup1)
	consumer := messaging.NewNotifConsumer(usecases.NotifUsecase)
	messaging.ConsumeTopic(ctx, consumerGroup, topic.Notif, consumer.ConsumeNotifEvent)
}
