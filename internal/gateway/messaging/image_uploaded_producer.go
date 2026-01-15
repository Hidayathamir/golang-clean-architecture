package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/consttopic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/ImageUploadedProducer.go -pkg=mock . ImageUploadedProducer

type ImageUploadedProducer interface {
	Send(ctx context.Context, event *model.ImageUploadedEvent) error
}

var _ ImageUploadedProducer = &ImageUploadedProducerImpl{}

type ImageUploadedProducerImpl struct {
	Config   *viper.Viper
	Producer sarama.SyncProducer
	Topic    string
}

func NewImageUploadedProducer(cfg *viper.Viper, producer sarama.SyncProducer) *ImageUploadedProducerImpl {
	return &ImageUploadedProducerImpl{
		Config:   cfg,
		Producer: producer,
		Topic:    consttopic.ImageUploaded,
	}
}

func (p *ImageUploadedProducerImpl) Send(ctx context.Context, event *model.ImageUploadedEvent) error {
	if p.Producer == nil {
		x.Logger.Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	message := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.ByteEncoder(value),
	}

	telemetry.InjectCtxToProducerMessage(ctx, message)

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	x.Logger.Debugf("Message sent to topic %s, partition %d, offset %d", p.Topic, partition, offset)

	return nil
}
