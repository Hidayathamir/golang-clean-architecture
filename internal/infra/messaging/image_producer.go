package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/twmb/franz-go/pkg/kgo"
)

//go:generate moq -out=../../mock/MockProducerImage2.go -pkg=mock . ImageProducer

type ImageProducer interface {
	SendImageUploaded(ctx context.Context, event *dto.ImageUploadedEvent) error
	SendImageLiked(ctx context.Context, event *dto.ImageLikedEvent) error
	SendImageCommented(ctx context.Context, event *dto.ImageCommentedEvent) error
}

var _ ImageProducer = &ImageProducerImpl{}

type ImageProducerImpl struct {
	Cfg    *config.Config
	Client *kgo.Client
}

func NewImageProducer(cfg *config.Config, client *kgo.Client) *ImageProducerImpl {
	return &ImageProducerImpl{
		Cfg:    cfg,
		Client: client,
	}
}

func (p *ImageProducerImpl) SendImageUploaded(ctx context.Context, event *dto.ImageUploadedEvent) error {
	err := p.send(ctx, topic.ImageUploaded, event)
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (p *ImageProducerImpl) SendImageLiked(ctx context.Context, event *dto.ImageLikedEvent) error {
	err := p.send(ctx, topic.ImageLiked, event)
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (p *ImageProducerImpl) SendImageCommented(ctx context.Context, event *dto.ImageCommentedEvent) error {
	err := p.send(ctx, topic.ImageCommented, event)
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (p *ImageProducerImpl) send(ctx context.Context, topicName string, event any) error {
	if p.Client == nil {
		x.Logger.WithContext(ctx).Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	record := &kgo.Record{
		Topic: topicName,
		Value: value,
	}

	if err := p.Client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return errkit.AddFuncName(err)
	}

	x.Logger.WithContext(ctx).WithField("topic", topicName).Debug("message sent")

	return nil
}
