package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/TodoProducer.go -pkg=mock . TodoProducer

type TodoProducer interface {
	Send(ctx context.Context, event *model.TodoCompletedEvent) error
}

var _ TodoProducer = &TodoProducerImpl{}

type TodoProducerImpl struct {
	Config   *viper.Viper
	Log      *logrus.Logger
	Producer sarama.SyncProducer
	Topic    string
}

func NewTodoProducer(cfg *viper.Viper, log *logrus.Logger, producer sarama.SyncProducer) *TodoProducerImpl {
	return &TodoProducerImpl{
		Config:   cfg,
		Log:      log,
		Producer: producer,
		Topic:    "todos",
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
		Key:   sarama.StringEncoder(event.GetID()),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName("messaging.(*TodoProducerImpl).Send", err)
	}

	p.Log.Debugf("Message sent to topic %s, partition %d, offset %d", p.Topic, partition, offset)
	return nil
}
