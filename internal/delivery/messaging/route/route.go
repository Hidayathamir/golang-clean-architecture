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

func ConsumeUserFollowedEventForNotification(ctx context.Context, cfg *config.Config, consumers *dependency_injection.Consumers) {
	x.Logger.Info("setup consumer.ConsumeUserFollowedEventForNotification")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.UserFollowedGroup1)
	messaging.ConsumeTopic(ctx, consumerGroup, topic.UserFollowed, consumers.UserConsumer.ConsumeUserFollowedEventForNotification)
}

func ConsumeUserFollowedEventForUpdateCount(ctx context.Context, cfg *config.Config, consumers *dependency_injection.Consumers) {
	x.Logger.Info("setup consumer.ConsumeUserFollowedEventForUpdateCount")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.UserFollowedGroup2)
	messaging.ConsumeTopicBatch(ctx, consumerGroup, topic.UserFollowed, consumers.UserConsumer.ConsumeUserFollowedEventForUpdateCount, 10, 10*time.Second)
}

func SetupImageUploadedConsumer(ctx context.Context, cfg *config.Config, consumers *dependency_injection.Consumers) {
	x.Logger.Info("setup consumer.ConsumeImageUploadedEvent")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.ImageUploadedGroup1)
	messaging.ConsumeTopic(ctx, consumerGroup, topic.ImageUploaded, consumers.ImageConsumer.ConsumeImageUploadedEvent)
}

func ConsumeImageLikedEventForNotification(ctx context.Context, cfg *config.Config, consumers *dependency_injection.Consumers) {
	x.Logger.Info("setup consumer.ConsumeImageLikedEventForNotification")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.ImageLikedGroup1)
	messaging.ConsumeTopic(ctx, consumerGroup, topic.ImageLiked, consumers.ImageConsumer.ConsumeImageLikedEventForNotification)
}

func ConsumeImageLikedEventForUpdateCount(ctx context.Context, cfg *config.Config, consumers *dependency_injection.Consumers) {
	x.Logger.Info("setup consumer.ConsumeImageLikedEventForUpdateCount")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.ImageLikedGroup2)
	messaging.ConsumeTopicBatch(ctx, consumerGroup, topic.ImageLiked, consumers.ImageConsumer.ConsumeImageLikedEventForUpdateCount, 10, 10*time.Second)
}

func ConsumeImageCommentedEventForNotification(ctx context.Context, cfg *config.Config, consumers *dependency_injection.Consumers) {
	x.Logger.Info("setup consumer.ConsumeImageCommentedEventForNotification")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.ImageCommentedGroup1)
	messaging.ConsumeTopic(ctx, consumerGroup, topic.ImageCommented, consumers.ImageConsumer.ConsumeImageCommentedEventForNotification)
}

func ConsumeImageCommentedEventForUpdateCount(ctx context.Context, cfg *config.Config, consumers *dependency_injection.Consumers) {
	x.Logger.Info("setup consumer.ConsumeImageCommentedEventForUpdateCount")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.ImageCommentedGroup2)
	messaging.ConsumeTopicBatch(ctx, consumerGroup, topic.ImageCommented, consumers.ImageConsumer.ConsumeImageCommentedEventForUpdateCount, 10, 18*time.Second)
}

func SetupNotifConsumer(ctx context.Context, cfg *config.Config, consumers *dependency_injection.Consumers) {
	x.Logger.Info("setup consumer.ConsumeNotifEvent")
	consumerGroup := provider.NewKafkaConsumerGroup(cfg, consumergroup.NotifGroup1)
	messaging.ConsumeTopic(ctx, consumerGroup, topic.Notif, consumers.NotifConsumer.ConsumeNotifEvent)
}
