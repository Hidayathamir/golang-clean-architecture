package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/image"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
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

	event := new(model.ImageUploadedEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
		return errkit.AddFuncName(err)
	}

	req := new(model.NotifyFollowerOnUploadRequest)
	converter.ModelImageUploadedEventToModelNotifyFollowerOnUploadRequest(ctx, event, req)

	if err := c.Usecase.NotifyFollowerOnUpload(ctx, req); err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *ImageConsumer) ConsumeImageLikedEvent(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	event := new(model.ImageLikedEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error("")
		return errkit.AddFuncName(err)
	}

	// TODO process event
	x.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"event":     event,
		"partition": message.Partition,
	}).Info("")

	return nil
}

func (c *ImageConsumer) ConsumeImageCommentedEvent(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	event := new(model.ImageCommentedEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error("")
		return errkit.AddFuncName(err)
	}

	// TODO process event
	x.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"event":     event,
		"partition": message.Partition,
	}).Info("")

	return nil
}
