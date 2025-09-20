package messaging

import (
	"encoding/json"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/errkit"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

//go:generate moq -out=../../mock/Producer.go -pkg=mock . Producer

type Producer[T model.Event] interface {
	GetTopic() *string
	Send(event T) error
}

var _ Producer[model.Event] = &ProducerImpl[model.Event]{}

type ProducerImpl[T model.Event] struct {
	Producer sarama.SyncProducer
	Topic    string
	Log      *logrus.Logger
}

func (p *ProducerImpl[T]) GetTopic() *string {
	return &p.Topic
}

func (p *ProducerImpl[T]) Send(event T) error {
	if p.Producer == nil {
		p.Log.Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		p.Log.WithError(err).Error("failed to marshal event")
		return errkit.AddFuncName(err)
	}

	message := &sarama.ProducerMessage{
		Topic: p.Topic,
		Key:   sarama.StringEncoder(event.GetId()),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		p.Log.WithError(err).Error("failed to produce message")
		return errkit.AddFuncName(err)
	}

	p.Log.Debugf("Message sent to topic %s, partition %d, offset %d", p.Topic, partition, offset)
	return nil
}
