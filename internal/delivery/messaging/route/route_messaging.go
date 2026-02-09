package route

import (
	"context"
	"sync"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/consumergroup"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
)

func Setup(ctx context.Context, cfg *config.Config, consumers *dependency_injection.Consumers, wg *sync.WaitGroup) {
	wg.Go(func() {
		consumerGroup := consumergroup.UserFollowedGroup1
		_topic := topic.UserFollowed
		handler := consumers.UserConsumer.NotifyUserBeingFollowed
		messaging.ConsumeEvent(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.UserFollowedGroup2
		_topic := topic.UserFollowed
		handler := consumers.UserConsumer.BatchUpdateUserFollowStats
		messaging.ConsumeEvent(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageUploadedGroup1
		_topic := topic.ImageUploaded
		handler := consumers.ImageConsumer.NotifyFollowerOnUpload
		messaging.ConsumeEvent(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageLikedGroup1
		_topic := topic.ImageLiked
		handler := consumers.ImageConsumer.NotifyUserImageLiked
		messaging.ConsumeEvent(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageLikedGroup2
		_topic := topic.ImageLiked
		handler := consumers.ImageConsumer.BatchUpdateImageLikeCount
		messaging.ConsumeEvent(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageCommentedGroup1
		_topic := topic.ImageCommented
		handler := consumers.ImageConsumer.NotifyUserImageCommented
		messaging.ConsumeEvent(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageCommentedGroup2
		_topic := topic.ImageCommented
		handler := consumers.ImageConsumer.BatchUpdateImageCommentCount
		messaging.ConsumeEvent(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.NotifGroup1
		_topic := topic.Notif
		handler := consumers.NotifConsumer.Notify
		messaging.ConsumeEvent(ctx, cfg, consumerGroup, _topic, handler)
	})
}
