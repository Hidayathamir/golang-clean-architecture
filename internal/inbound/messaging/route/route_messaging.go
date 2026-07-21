package route

import (
	"context"
	"sync"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/consumergroup"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
)

func Setup(ctx context.Context, cfg *config.Config, consumers *dependency_injection.Consumers, wg *sync.WaitGroup) {
	wg.Go(func() {
		consumerGroup := consumergroup.UserFollowedNotifyUser
		_topic := topic.UserFollowed
		handler := consumers.UserConsumer.NotifyUserBeingFollowed
		messaging.ConsumeEventSingle(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.UserFollowedBatchStats
		_topic := topic.UserFollowed
		handler := consumers.UserConsumer.BatchUpdateUserFollowStats
		messaging.ConsumeEventBatch(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageUploadedNotifyFollowers
		_topic := topic.ImageUploaded
		handler := consumers.ImageConsumer.NotifyFollowerOnUpload
		messaging.ConsumeEventSingle(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageUploadedSyncSearch
		_topic := topic.ImageUploaded
		handler := consumers.ImageConsumer.SyncImageToElasticsearch
		messaging.ConsumeEventSingle(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageLikedNotifyOwner
		_topic := topic.ImageLiked
		handler := consumers.ImageConsumer.NotifyUserImageLiked
		messaging.ConsumeEventSingle(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageLikedBatchCount
		_topic := topic.ImageLiked
		handler := consumers.ImageConsumer.BatchUpdateImageLikeCount
		messaging.ConsumeEventBatch(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageCommentedNotifyOwner
		_topic := topic.ImageCommented
		handler := consumers.ImageConsumer.NotifyUserImageCommented
		messaging.ConsumeEventSingle(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageCommentedBatchCount
		_topic := topic.ImageCommented
		handler := consumers.ImageConsumer.BatchUpdateImageCommentCount
		messaging.ConsumeEventBatch(ctx, cfg, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.NotifLog
		_topic := topic.Notif
		handler := consumers.NotifConsumer.Notify
		messaging.ConsumeEventSingle(ctx, cfg, consumerGroup, _topic, handler)
	})
}
