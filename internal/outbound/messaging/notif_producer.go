package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/repository"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockProducerNotif2.go -pkg=mock . NotifProducer

type NotifProducer interface {
	SendNotif(ctx context.Context, db *gorm.DB, event *dto.NotifEvent) error
}

var _ NotifProducer = &NotifProducerImpl{}

type NotifProducerImpl struct {
	Cfg              *config.Config
	OutboxRepository repository.OutboxRepository
}

func NewNotifProducer(cfg *config.Config, outboxRepository repository.OutboxRepository) *NotifProducerImpl {
	return &NotifProducerImpl{
		Cfg:              cfg,
		OutboxRepository: outboxRepository,
	}
}

func (p *NotifProducerImpl) SendNotif(ctx context.Context, db *gorm.DB, event *dto.NotifEvent) error {
	err := p.send(ctx, db, topic.Notif, event)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*NotifProducerImpl).SendNotif")
	}
	return nil
}

func (p *NotifProducerImpl) send(ctx context.Context, db *gorm.DB, topicName string, event any) error {
	if !p.Cfg.GetKafkaProducerEnabled() {
		logkit.Logger.WithContext(ctx).Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*NotifProducerImpl).send")
	}

	outbox := entity.Outbox{
		Topic:        topicName,
		Payload:      value,
		TraceContext: telemetry.InjectTraceContext(ctx),
		Status:       entity.OutboxStatusPending,
	}

	err = p.OutboxRepository.Insert(ctx, db, &outbox)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*NotifProducerImpl).send")
	}

	logkit.Logger.WithContext(ctx).WithField("topic", topicName).Debug("outbox record inserted")

	return nil
}
