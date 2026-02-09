package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
)

//go:generate moq -out=../../mock/MockProducerUser.go -pkg=mock . UserProducer

type UserProducer interface {
	SendUserFollowed(ctx context.Context, event *dto.UserFollowedEvent) error
}

var _ UserProducer = &UserProducerImpl{}

type UserProducerImpl struct {
	Cfg      *config.Config
	Producer sarama.SyncProducer
}

func NewUserProducer(cfg *config.Config, producer sarama.SyncProducer) *UserProducerImpl {
	return &UserProducerImpl{
		Cfg:      cfg,
		Producer: producer,
	}
}

func (p *UserProducerImpl) SendUserFollowed(ctx context.Context, event *dto.UserFollowedEvent) error {
	if p.Producer == nil {
		x.Logger.WithContext(ctx).Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	message := &sarama.ProducerMessage{
		Topic: topic.UserFollowed,
		Key:   sarama.StringEncoder(event.GetID()),
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
