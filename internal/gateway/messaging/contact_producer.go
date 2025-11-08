package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/ContactProducer.go -pkg=mock . ContactProducer

type ContactProducer interface {
	Send(ctx context.Context, event *model.ContactEvent) error
}

var _ ContactProducer = &ContactProducerImpl{}

type ContactProducerImpl struct {
	Config   *viper.Viper
	Producer sarama.SyncProducer
	Topic    string
}

func NewContactProducer(cfg *viper.Viper, producer sarama.SyncProducer) *ContactProducerImpl {
	return &ContactProducerImpl{
		Config:   cfg,
		Producer: producer,
		Topic:    "contacts",
	}
}

func (p *ContactProducerImpl) Send(ctx context.Context, event *model.ContactEvent) error {
	if p.Producer == nil {
		x.Logger.Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName("messaging.(*ContactProducerImpl).Send", err)
	}

	message := &sarama.ProducerMessage{
		Topic: p.Topic,
		Key:   sarama.StringEncoder(event.GetID()),
		Value: sarama.ByteEncoder(value),
	}

	telemetry.InjectCtxToProducerMessage(ctx, message)

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName("messaging.(*ContactProducerImpl).Send", err)
	}

	x.Logger.Debugf("Message sent to topic %s, partition %d, offset %d", p.Topic, partition, offset)
	return nil
}
