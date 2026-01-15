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

//go:generate moq -out=../../mock/ImageCommentedProducer.go -pkg=mock . ImageCommentedProducer

type ImageCommentedProducer interface {
	Send(ctx context.Context, event *model.ImageCommentedEvent) error
}

var _ ImageCommentedProducer = &ImageCommentedProducerImpl{}

type ImageCommentedProducerImpl struct {
	Config   *viper.Viper
	Producer sarama.SyncProducer
	Topic    string
}

func NewImageCommentedProducer(cfg *viper.Viper, producer sarama.SyncProducer) *ImageCommentedProducerImpl {
	return &ImageCommentedProducerImpl{
		Config:   cfg,
		Producer: producer,
		Topic:    consttopic.ImageCommented,
	}
}

func (p *ImageCommentedProducerImpl) Send(ctx context.Context, event *model.ImageCommentedEvent) error {
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
