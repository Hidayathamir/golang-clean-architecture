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

//go:generate moq -out=../../mock/ImageLikedProducer.go -pkg=mock . ImageLikedProducer

type ImageLikedProducer interface {
	Send(ctx context.Context, event *model.ImageLikedEvent) error
}

var _ ImageLikedProducer = &ImageLikedProducerImpl{}

type ImageLikedProducerImpl struct {
	Config   *viper.Viper
	Producer sarama.SyncProducer
	Topic    string
}

func NewImageLikedProducer(cfg *viper.Viper, producer sarama.SyncProducer) *ImageLikedProducerImpl {
	return &ImageLikedProducerImpl{
		Config:   cfg,
		Producer: producer,
		Topic:    consttopic.ImageLiked,
	}
}

func (p *ImageLikedProducerImpl) Send(ctx context.Context, event *model.ImageLikedEvent) error {
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
