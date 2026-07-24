package route

import (
	"context"
	"sync"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dependency_injection"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/consumergroup"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
	"github.com/twmb/franz-go/pkg/kgo"
)

func Setup(ctx context.Context, cfg *config.Config, producer *kgo.Client, consumers *dependency_injection.Consumers, wg *sync.WaitGroup) {
	// --- primary consumers: single ---

	wg.Go(func() {
		consumerGroup := consumergroup.UserFollowedNotifyUser
		_topic := topic.UserFollowed
		handler := consumers.UserConsumer.NotifyUserBeingFollowed
		messaging.ConsumeEventSingle(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageUploadedNotifyFollowers
		_topic := topic.ImageUploaded
		handler := consumers.ImageConsumer.NotifyFollowerOnUpload
		messaging.ConsumeEventSingle(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageUploadedSyncSearch
		_topic := topic.ImageUploaded
		handler := consumers.ImageConsumer.SyncImageToElasticsearch
		messaging.ConsumeEventSingle(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageLikedNotifyOwner
		_topic := topic.ImageLiked
		handler := consumers.ImageConsumer.NotifyUserImageLiked
		messaging.ConsumeEventSingle(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageCommentedNotifyOwner
		_topic := topic.ImageCommented
		handler := consumers.ImageConsumer.NotifyUserImageCommented
		messaging.ConsumeEventSingle(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.NotifLog
		_topic := topic.Notif
		handler := consumers.NotifConsumer.Notify
		messaging.ConsumeEventSingle(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	// --- primary consumers: batch ---

	wg.Go(func() {
		consumerGroup := consumergroup.UserFollowedBatchStats
		_topic := topic.UserFollowed
		handler := consumers.UserConsumer.BatchUpdateUserFollowStats
		messaging.ConsumeEventBatch(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageLikedBatchCount
		_topic := topic.ImageLiked
		handler := consumers.ImageConsumer.BatchUpdateImageLikeCount
		messaging.ConsumeEventBatch(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageCommentedBatchCount
		_topic := topic.ImageCommented
		handler := consumers.ImageConsumer.BatchUpdateImageCommentCount
		messaging.ConsumeEventBatch(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	// --- retry consumers: single handlers ---

	wg.Go(func() {
		consumerGroup := consumergroup.UserFollowedNotifyUserRetry
		_topic := topic.UserFollowed
		handler := consumers.UserConsumer.NotifyUserBeingFollowed
		messaging.ConsumeEventRetry(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageUploadedNotifyFollowersRetry
		_topic := topic.ImageUploaded
		handler := consumers.ImageConsumer.NotifyFollowerOnUpload
		messaging.ConsumeEventRetry(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageUploadedSyncSearchRetry
		_topic := topic.ImageUploaded
		handler := consumers.ImageConsumer.SyncImageToElasticsearch
		messaging.ConsumeEventRetry(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageLikedNotifyOwnerRetry
		_topic := topic.ImageLiked
		handler := consumers.ImageConsumer.NotifyUserImageLiked
		messaging.ConsumeEventRetry(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageCommentedNotifyOwnerRetry
		_topic := topic.ImageCommented
		handler := consumers.ImageConsumer.NotifyUserImageCommented
		messaging.ConsumeEventRetry(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.NotifLogRetry
		_topic := topic.Notif
		handler := consumers.NotifConsumer.Notify
		messaging.ConsumeEventRetry(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	// --- retry consumers: batch handlers ---

	wg.Go(func() {
		consumerGroup := consumergroup.UserFollowedBatchStatsRetry
		_topic := topic.UserFollowed
		handler := consumers.UserConsumer.UpdateUserFollowStats
		messaging.ConsumeEventRetry(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageLikedBatchCountRetry
		_topic := topic.ImageLiked
		handler := consumers.ImageConsumer.UpdateImageLikeCount
		messaging.ConsumeEventRetry(ctx, cfg, producer, consumerGroup, _topic, handler)
	})

	wg.Go(func() {
		consumerGroup := consumergroup.ImageCommentedBatchCountRetry
		_topic := topic.ImageCommented
		handler := consumers.ImageConsumer.UpdateImageCommentCount
		messaging.ConsumeEventRetry(ctx, cfg, producer, consumerGroup, _topic, handler)
	})
}
