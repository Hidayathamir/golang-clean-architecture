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

//go:generate moq -out=../../mock/UserProducer.go -pkg=mock . UserProducer

type UserProducer interface {
	Send(ctx context.Context, event *model.UserEvent) error
}

var _ UserProducer = &UserProducerImpl{}

type UserProducerImpl struct {
	Config   *viper.Viper
	Producer sarama.SyncProducer
	Topic    string
}

func NewUserProducer(cfg *viper.Viper, producer sarama.SyncProducer) *UserProducerImpl {
	return &UserProducerImpl{
		Config:   cfg,
		Producer: producer,
		Topic:    "users",
	}
}

func (p *UserProducerImpl) Send(ctx context.Context, event *model.UserEvent) error {
	if p.Producer == nil {
		x.Logger.Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName("messaging.(*UserProducerImpl).Send", err)
	}

	message := &sarama.ProducerMessage{
		Topic: p.Topic,
		Key:   sarama.StringEncoder(event.GetID()),
		Value: sarama.ByteEncoder(value),
	}

	telemetry.InjectCtxToProducerMessage(ctx, message)

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName("messaging.(*UserProducerImpl).Send", err)
	}

	x.Logger.Debugf("Message sent to topic %s, partition %d, offset %d", p.Topic, partition, offset)
	return nil
}
