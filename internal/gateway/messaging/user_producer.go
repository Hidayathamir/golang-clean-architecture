package messaging

import (
	"context"
	"encoding/json"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/errkit"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

//go:generate moq -out=../../mock/UserProducer.go -pkg=mock . UserProducer

type UserProducer interface {
	Send(ctx context.Context, event *model.UserEvent) error
}

type UserProducerImpl struct {
	Producer sarama.SyncProducer
	Topic    string
	Log      *logrus.Logger
}

func NewUserProducer(producer sarama.SyncProducer, log *logrus.Logger) *UserProducerImpl {
	return &UserProducerImpl{
		Producer: producer,
		Topic:    "users",
		Log:      log,
	}
}

func (p *UserProducerImpl) Send(ctx context.Context, event *model.UserEvent) error {
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
