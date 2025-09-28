package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

//go:generate moq -out=../../mock/ContactProducer.go -pkg=mock . ContactProducer

type ContactProducer interface {
	Send(ctx context.Context, event *model.ContactEvent) error
}

var _ ContactProducer = &ContactProducerImpl{}

type ContactProducerImpl struct {
	Producer sarama.SyncProducer
	Topic    string
	Log      *logrus.Logger
}

func NewContactProducer(producer sarama.SyncProducer, log *logrus.Logger) *ContactProducerImpl {
	return &ContactProducerImpl{
		Producer: producer,
		Topic:    "contacts",
		Log:      log,
	}
}

func (p *ContactProducerImpl) Send(ctx context.Context, event *model.ContactEvent) error {
	if p.Producer == nil {
		p.Log.Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	message := &sarama.ProducerMessage{
		Topic: p.Topic,
		Key:   sarama.StringEncoder(event.GetId()),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	p.Log.Debugf("Message sent to topic %s, partition %d, offset %d", p.Topic, partition, offset)
	return nil
}
