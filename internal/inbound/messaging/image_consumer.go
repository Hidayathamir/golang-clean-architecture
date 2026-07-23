package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/imageusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ImageConsumer struct {
	Usecase imageusecase.ImageUsecase
}

func NewImageConsumer(usecase imageusecase.ImageUsecase) *ImageConsumer {
	return &ImageConsumer{
		Usecase: usecase,
	}
}

func (c *ImageConsumer) NotifyFollowerOnUpload(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.ImageUploadedEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(errkit.WrapNonRetryable(err), "messaging.(*ImageConsumer).NotifyFollowerOnUpload")
	}

	req := dto.NotifyFollowerOnUploadRequest{}
	converter.DtoImageUploadedEventToDtoNotifyFollowerOnUploadRequest(event, &req)

	err = c.Usecase.NotifyFollowerOnUpload(ctx, req)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*ImageConsumer).NotifyFollowerOnUpload")
	}

	return nil
}

func (c *ImageConsumer) SyncImageToElasticsearch(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.ImageUploadedEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(errkit.WrapNonRetryable(err), "messaging.(*ImageConsumer).SyncImageToElasticsearch")
	}

	req := dto.SyncImageToElasticsearchRequest{}
	converter.DtoImageUploadedEventToDtoSyncImageToElasticsearchRequest(event, &req)

	err = c.Usecase.SyncImageToElasticsearch(ctx, req)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*ImageConsumer).SyncImageToElasticsearch")
	}

	return nil
}

func (c *ImageConsumer) NotifyUserImageLiked(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.ImageLikedEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(errkit.WrapNonRetryable(err), "messaging.(*ImageConsumer).NotifyUserImageLiked")
	}

	req := dto.NotifyUserImageLikedRequest{}
	converter.DtoImageLikedEventToDtoNotifyUserImageLikedRequest(event, &req)

	err = c.Usecase.NotifyUserImageLiked(ctx, req)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*ImageConsumer).NotifyUserImageLiked")
	}

	return nil
}

func (c *ImageConsumer) BatchUpdateImageLikeCount(originalCtx context.Context, records []*kgo.Record) error {
	ctx, span := telemetry.StartConsumerBatch(originalCtx, records)
	defer span.End()

	req := dto.BatchUpdateImageLikeCountRequest{}
	converter.KGoRecordListToDtoBatchUpdateImageLikeCountRequest(ctx, records, &req)

	err := c.Usecase.BatchUpdateImageLikeCount(ctx, req)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*ImageConsumer).BatchUpdateImageLikeCount")
	}

	return nil
}

func (c *ImageConsumer) UpdateImageLikeCount(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	err := c.BatchUpdateImageLikeCount(ctx, []*kgo.Record{record})
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*ImageConsumer).UpdateImageLikeCount")
	}

	return nil
}

func (c *ImageConsumer) NotifyUserImageCommented(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.ImageCommentedEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(errkit.WrapNonRetryable(err), "messaging.(*ImageConsumer).NotifyUserImageCommented")
	}

	req := dto.NotifyUserImageCommentedRequest{}
	converter.DtoImageCommentedEventToDtoNotifyUserImageCommentedRequest(event, &req)

	err = c.Usecase.NotifyUserImageCommented(ctx, req)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*ImageConsumer).NotifyUserImageCommented")
	}

	return nil
}

func (c *ImageConsumer) BatchUpdateImageCommentCount(ctx context.Context, records []*kgo.Record) error {
	ctx, span := telemetry.StartConsumerBatch(ctx, records)
	defer span.End()

	req := dto.BatchUpdateImageCommentCountRequest{}
	converter.KGoRecordListToDtoBatchUpdateImageCommentCountRequest(ctx, records, &req)

	err := c.Usecase.BatchUpdateImageCommentCount(ctx, req)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*ImageConsumer).BatchUpdateImageCommentCount")
	}

	return nil
}

func (c *ImageConsumer) UpdateImageCommentCount(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	err := c.BatchUpdateImageCommentCount(ctx, []*kgo.Record{record})
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*ImageConsumer).UpdateImageCommentCount")
	}

	return nil
}
