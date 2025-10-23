package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

//go:generate moq -out=../../mock/TodoProducer.go -pkg=mock . TodoProducer

type TodoProducer interface {
	Send(ctx context.Context, event *model.TodoCompletedEvent) error
}

var _ TodoProducer = &TodoProducerImpl{}

type TodoProducerImpl struct {
	Producer sarama.SyncProducer
	Topic    string
	Log      *logrus.Logger
}

func NewTodoProducer(producer sarama.SyncProducer, log *logrus.Logger) *TodoProducerImpl {
	return &TodoProducerImpl{
		Producer: producer,
		Topic:    "todos",
		Log:      log,
	}
}

func (p *TodoProducerImpl) Send(ctx context.Context, event *model.TodoCompletedEvent) error {
	if p.Producer == nil {
		p.Log.Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName("messaging.(*TodoProducerImpl).Send", err)
	}

	message := &sarama.ProducerMessage{
		Topic: p.Topic,
		Key:   sarama.StringEncoder(event.GetId()),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName("messaging.(*TodoProducerImpl).Send", err)
	}

	p.Log.Debugf("Message sent to topic %s, partition %d, offset %d", p.Topic, partition, offset)
	return nil
}
