package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/twmb/franz-go/pkg/kgo"
)

//go:generate moq -out=../../mock/MockProducerNotif2.go -pkg=mock . NotifProducer

type NotifProducer interface {
	SendNotif(ctx context.Context, event *dto.NotifEvent) error
}

var _ NotifProducer = &NotifProducerImpl{}

type NotifProducerImpl struct {
	Cfg    *config.Config
	Client *kgo.Client
}

func NewNotifProducer(cfg *config.Config, client *kgo.Client) *NotifProducerImpl {
	return &NotifProducerImpl{
		Cfg:    cfg,
		Client: client,
	}
}

func (p *NotifProducerImpl) SendNotif(ctx context.Context, event *dto.NotifEvent) error {
	err := p.send(ctx, topic.Notif, event)
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}

func (p *NotifProducerImpl) send(ctx context.Context, topicName string, event any) error {
	if p.Client == nil {
		x.Logger.WithContext(ctx).Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	record := &kgo.Record{
		Topic: topicName,
		Value: value,
	}

	err = p.Client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		return errkit.AddFuncName(err)
	}

	x.Logger.WithContext(ctx).WithField("topic", topicName).Debug("message sent")

	return nil
}
