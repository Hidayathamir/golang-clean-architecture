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
		if err := c.notifyFollowerOnUpload(ctx, record); err != nil {
			x.Logger.WithContext(ctx).WithError(err).Error()
			continue
		}
	}
	return nil
}

func (c *ImageConsumer) notifyFollowerOnUpload(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := new(dto.ImageUploadedEvent)
	if err := json.Unmarshal(record.Value, event); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := new(dto.NotifyFollowerOnUploadRequest)
	converter.DtoImageUploadedEventToDtoNotifyFollowerOnUploadRequest(ctx, event, req)

	if err := c.Usecase.NotifyFollowerOnUpload(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *ImageConsumer) SyncImageToElasticsearch(ctx context.Context, records []*kgo.Record) error {
	for _, record := range records {
		if err := c.syncImageToElasticsearch(ctx, record); err != nil {
			x.Logger.WithContext(ctx).WithError(err).Error()
			continue
		}
	}
	return nil
}

func (c *ImageConsumer) syncImageToElasticsearch(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := new(dto.ImageUploadedEvent)
	if err := json.Unmarshal(record.Value, event); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := new(dto.SyncImageToElasticsearchRequest)
	converter.DtoImageUploadedEventToDtoSyncImageToElasticsearchRequest(ctx, event, req)

	if err := c.Usecase.SyncImageToElasticsearch(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *ImageConsumer) NotifyUserImageLiked(ctx context.Context, records []*kgo.Record) error {
	for _, record := range records {
		if err := c.notifyUserImageLiked(ctx, record); err != nil {
			x.Logger.WithContext(ctx).WithError(err).Error()
			continue
		}
	}
	return nil
}

func (c *ImageConsumer) notifyUserImageLiked(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := new(dto.ImageLikedEvent)
	if err := json.Unmarshal(record.Value, event); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := new(dto.NotifyUserImageLikedRequest)
	converter.DtoImageLikedEventToDtoNotifyUserImageLikedRequest(ctx, event, req)

	if err := c.Usecase.NotifyUserImageLiked(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *ImageConsumer) BatchUpdateImageLikeCount(originalCtx context.Context, records []*kgo.Record) error {
	ctx, span := telemetry.StartConsumerBatch(originalCtx, records)
	defer span.End()

	req := new(dto.BatchUpdateImageLikeCountRequest)
	converter.KGoRecordListToDtoBatchUpdateImageLikeCountRequest(ctx, records, req)

	if err := c.Usecase.BatchUpdateImageLikeCount(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *ImageConsumer) NotifyUserImageCommented(ctx context.Context, records []*kgo.Record) error {
	for _, record := range records {
		if err := c.notifyUserImageCommented(ctx, record); err != nil {
			x.Logger.WithContext(ctx).WithError(err).Error()
			continue
		}
	}
	return nil
}

func (c *ImageConsumer) notifyUserImageCommented(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := new(dto.ImageCommentedEvent)
	if err := json.Unmarshal(record.Value, event); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := new(dto.NotifyUserImageCommentedRequest)
	converter.DtoImageCommentedEventToDtoNotifyUserImageCommentedRequest(ctx, event, req)

	if err := c.Usecase.NotifyUserImageCommented(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *ImageConsumer) BatchUpdateImageCommentCount(ctx context.Context, records []*kgo.Record) error {
	ctx, span := telemetry.StartConsumerBatch(ctx, records)
	defer span.End()

	req := new(dto.BatchUpdateImageCommentCountRequest)
	converter.KGoRecordListToDtoBatchUpdateImageCommentCountRequest(ctx, records, req)

	if err := c.Usecase.BatchUpdateImageCommentCount(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}
