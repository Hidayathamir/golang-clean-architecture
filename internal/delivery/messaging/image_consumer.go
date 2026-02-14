package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/image"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ImageConsumer struct {
	Usecase image.ImageUsecase
}

func NewImageConsumer(usecase image.ImageUsecase) *ImageConsumer {
	return &ImageConsumer{
		Usecase: usecase,
	}
}

func (c *ImageConsumer) NotifyFollowerOnUpload(ctx context.Context, records []*kgo.Record) error {
	for _, record := range records {
		err := c.notifyFollowerOnUpload(ctx, record)
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).Error()
			continue
		}
	}
	return nil
}

func (c *ImageConsumer) notifyFollowerOnUpload(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.ImageUploadedEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := dto.NotifyFollowerOnUploadRequest{}
	converter.DtoImageUploadedEventToDtoNotifyFollowerOnUploadRequest(event, &req)

	err = c.Usecase.NotifyFollowerOnUpload(ctx, req)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *ImageConsumer) SyncImageToElasticsearch(ctx context.Context, records []*kgo.Record) error {
	for _, record := range records {
		err := c.syncImageToElasticsearch(ctx, record)
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).Error()
			continue
		}
	}
	return nil
}

func (c *ImageConsumer) syncImageToElasticsearch(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.ImageUploadedEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := dto.SyncImageToElasticsearchRequest{}
	converter.DtoImageUploadedEventToDtoSyncImageToElasticsearchRequest(event, &req)

	err = c.Usecase.SyncImageToElasticsearch(ctx, req)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *ImageConsumer) NotifyUserImageLiked(ctx context.Context, records []*kgo.Record) error {
	for _, record := range records {
		err := c.notifyUserImageLiked(ctx, record)
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).Error()
			continue
		}
	}
	return nil
}

func (c *ImageConsumer) notifyUserImageLiked(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.ImageLikedEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := dto.NotifyUserImageLikedRequest{}
	converter.DtoImageLikedEventToDtoNotifyUserImageLikedRequest(event, &req)

	err = c.Usecase.NotifyUserImageLiked(ctx, req)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
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
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *ImageConsumer) NotifyUserImageCommented(ctx context.Context, records []*kgo.Record) error {
	for _, record := range records {
		err := c.notifyUserImageCommented(ctx, record)
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).Error()
			continue
		}
	}
	return nil
}

func (c *ImageConsumer) notifyUserImageCommented(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.ImageCommentedEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := dto.NotifyUserImageCommentedRequest{}
	converter.DtoImageCommentedEventToDtoNotifyUserImageCommentedRequest(event, &req)

	err = c.Usecase.NotifyUserImageCommented(ctx, req)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
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
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}
