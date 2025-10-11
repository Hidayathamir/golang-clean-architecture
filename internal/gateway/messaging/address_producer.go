package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

//go:generate moq -out=../../mock/AddressProducer.go -pkg=mock . AddressProducer

type AddressProducer interface {
	Send(ctx context.Context, event *model.AddressEvent) error
}

var _ AddressProducer = &AddressProducerImpl{}

type AddressProducerImpl struct {
	Producer sarama.SyncProducer
	Topic    string
	Log      *logrus.Logger
}

func NewAddressProducer(producer sarama.SyncProducer, log *logrus.Logger) *AddressProducerImpl {
	return &AddressProducerImpl{
		Producer: producer,
		Topic:    "addresses",
		Log:      log,
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
		Key:   sarama.StringEncoder(event.GetId()),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName("messaging.(*AddressProducerImpl).Send", err)
	}

	p.Log.Debugf("Message sent to topic %s, partition %d, offset %d", p.Topic, partition, offset)
	return nil
}
