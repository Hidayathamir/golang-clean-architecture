package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
)

//go:generate moq -out=../../mock/MockProducerImage.go -pkg=mock . ImageProducer

type ImageProducer interface {
	SendImageUploaded(ctx context.Context, event *model.ImageUploadedEvent) error
	SendImageLiked(ctx context.Context, event *model.ImageLikedEvent) error
	SendImageCommented(ctx context.Context, event *model.ImageCommentedEvent) error
}

var _ ImageProducer = &ImageProducerImpl{}

type ImageProducerImpl struct {
	Cfg      *config.Config
	Producer sarama.SyncProducer
}

func NewImageProducer(cfg *config.Config, producer sarama.SyncProducer) *ImageProducerImpl {
	return &ImageProducerImpl{
		Cfg:      cfg,
		Producer: producer,
	}
}

func (p *ImageProducerImpl) SendImageUploaded(ctx context.Context, event *model.ImageUploadedEvent) error {
	if p.Producer == nil {
		x.Logger.WithContext(ctx).Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic.ImageUploaded,
		Value: sarama.ByteEncoder(value),
	}

	telemetry.InjectCtxToProducerMessage(ctx, message)

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	x.Logger.WithContext(ctx).Debugf("Message sent to topic %s, partition %d, offset %d", message.Topic, partition, offset)

	return nil
}

func (p *ImageProducerImpl) SendImageLiked(ctx context.Context, event *model.ImageLikedEvent) error {
	if p.Producer == nil {
		x.Logger.WithContext(ctx).Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic.ImageLiked,
		Value: sarama.ByteEncoder(value),
	}

	telemetry.InjectCtxToProducerMessage(ctx, message)

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	x.Logger.WithContext(ctx).Debugf("Message sent to topic %s, partition %d, offset %d", message.Topic, partition, offset)

	return nil
}

func (p *ImageProducerImpl) SendImageCommented(ctx context.Context, event *model.ImageCommentedEvent) error {
	if p.Producer == nil {
		x.Logger.WithContext(ctx).Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic.ImageCommented,
		Value: sarama.ByteEncoder(value),
	}

	telemetry.InjectCtxToProducerMessage(ctx, message)

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	x.Logger.WithContext(ctx).Debugf("Message sent to topic %s, partition %d, offset %d", message.Topic, partition, offset)

	return nil
}
