package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/AddressProducer.go -pkg=mock . AddressProducer

type AddressProducer interface {
	Send(ctx context.Context, event *model.AddressEvent) error
}

var _ AddressProducer = &AddressProducerImpl{}

type AddressProducerImpl struct {
	Config   *viper.Viper
	Log      *logrus.Logger
	Producer sarama.SyncProducer
	Topic    string
}

func NewAddressProducer(cfg *viper.Viper, log *logrus.Logger, producer sarama.SyncProducer) *AddressProducerImpl {
	return &AddressProducerImpl{
		Config:   cfg,
		Log:      log,
		Producer: producer,
		Topic:    "addresses",
	}
}

func (p *AddressProducerImpl) Send(ctx context.Context, event *model.AddressEvent) error {
	if p.Producer == nil {
		p.Log.Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName("messaging.(*AddressProducerImpl).Send", err)
	}

	message := &sarama.ProducerMessage{
		Topic: p.Topic,
		Key:   sarama.StringEncoder(event.GetID()),
		Value: sarama.ByteEncoder(value),
	}

	telemetry.InjectCtxToProducerMessage(ctx, message)

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName("messaging.(*AddressProducerImpl).Send", err)
	}

	p.Log.Debugf("Message sent to topic %s, partition %d, offset %d", p.Topic, partition, offset)
	return nil
}
