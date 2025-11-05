package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/l"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/TodoProducer.go -pkg=mock . TodoProducer

type TodoProducer interface {
	Send(ctx context.Context, event *model.TodoCompletedEvent) error
}

var _ TodoProducer = &TodoProducerImpl{}

type TodoProducerImpl struct {
	Config   *viper.Viper
	Producer sarama.SyncProducer
	Topic    string
}

func NewTodoProducer(cfg *viper.Viper, producer sarama.SyncProducer) *TodoProducerImpl {
	return &TodoProducerImpl{
		Config:   cfg,
		Producer: producer,
		Topic:    "todos",
	}
}

func (p *TodoProducerImpl) Send(ctx context.Context, event *model.TodoCompletedEvent) error {
	if p.Producer == nil {
		l.Logger.Warn("Kafka producer is disabled")
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

	telemetry.InjectCtxToProducerMessage(ctx, message)

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName("messaging.(*TodoProducerImpl).Send", err)
	}

	l.Logger.Debugf("Message sent to topic %s, partition %d, offset %d", p.Topic, partition, offset)
	return nil
}
