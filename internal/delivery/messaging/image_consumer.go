package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/image"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
)

type ImageConsumer struct {
	Usecase image.ImageUsecase
}

func NewImageConsumer(usecase image.ImageUsecase) *ImageConsumer {
	return &ImageConsumer{
		Usecase: usecase,
	}
}

func (c *ImageConsumer) ConsumeImageUploadedEvent(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	event := new(dto.ImageUploadedEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
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

func (c *ImageConsumer) ConsumeImageLikedEventForNotification(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	event := new(dto.ImageLikedEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
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

func (c *ImageConsumer) ConsumeImageLikedEventForUpdateCount(messages []*sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartNew()
	defer span.End()

	req := new(dto.BatchUpdateImageLikeCountRequest)
	converter.SaramaConsumerMessageListToDtoBatchUpdateImageLikeCountRequest(ctx, messages, req)

	if err := c.Usecase.BatchUpdateImageLikeCount(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *ImageConsumer) ConsumeImageCommentedEventForNotification(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	event := new(dto.ImageCommentedEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
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

func (c *ImageConsumer) ConsumeImageCommentedEventForUpdateCount(messages []*sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartNew()
	defer span.End()

	req := new(dto.BatchUpdateImageCommentCountRequest)
	converter.SaramaConsumerMessageListToDtoBatchUpdateImageCommentCountRequest(ctx, messages, req)

	if err := c.Usecase.BatchUpdateImageCommentCount(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}
