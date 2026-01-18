package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/MockProducerNotif.go -pkg=mock . NotifProducer

type NotifProducer interface {
	SendNotif(ctx context.Context, event *model.NotifEvent) error
}

var _ NotifProducer = &NotifProducerImpl{}

type NotifProducerImpl struct {
	Config   *viper.Viper
	Producer sarama.SyncProducer
}

func NewNotifProducer(cfg *viper.Viper, producer sarama.SyncProducer) *NotifProducerImpl {
	return &NotifProducerImpl{
		Config:   cfg,
		Producer: producer,
	}
}

func (p *NotifProducerImpl) SendNotif(ctx context.Context, event *model.NotifEvent) error {
	if p.Producer == nil {
		x.Logger.WithContext(ctx).Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic.Notif,
		Value: sarama.ByteEncoder(value),
	}

	telemetry.InjectCtxToProducerMessage(ctx, message)

	partition, offset, err := p.Producer.SendMessage(message)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	x.Logger.WithContext(ctx).Debugf("Message sent to topic %s, partition %d, offset %d", message.Topic, partition, offset)

	return nil
}
