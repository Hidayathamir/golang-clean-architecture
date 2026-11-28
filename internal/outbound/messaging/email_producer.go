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

//go:generate moq -out=../../mock/MockProducerEmail.go -pkg=mock . EmailProducer

type EmailProducer interface {
	SendVerificationEmail(ctx context.Context, db *gorm.DB, event *dto.EmailVerificationEvent) error
}

var _ EmailProducer = &EmailProducerImpl{}

type EmailProducerImpl struct {
	Cfg              *config.Config
	OutboxRepository repository.OutboxRepository
}

func NewEmailProducer(cfg *config.Config, outboxRepository repository.OutboxRepository) *EmailProducerImpl {
	return &EmailProducerImpl{
		Cfg:              cfg,
		OutboxRepository: outboxRepository,
	}
}

func (p *EmailProducerImpl) SendVerificationEmail(ctx context.Context, db *gorm.DB, event *dto.EmailVerificationEvent) error {
	err := p.send(ctx, db, topic.EmailVerification, event)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*EmailProducerImpl).SendVerificationEmail")
	}
	return nil
}

func (p *EmailProducerImpl) send(ctx context.Context, db *gorm.DB, topicName topic.Topic, event any) error {
	if !p.Cfg.GetKafkaProducerEnabled() {
		logkit.Logger.WithContext(ctx).Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*EmailProducerImpl).send")
	}

	outbox := entity.Outbox{
		Topic:        topicName.Primary,
		Payload:      value,
		TraceContext: telemetry.InjectTraceContext(ctx),
		Status:       entity.OutboxStatusPending,
	}

	err = p.OutboxRepository.Insert(ctx, db, &outbox)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*EmailProducerImpl).send")
	}

	logkit.Logger.WithContext(ctx).WithField("topic", topicName.Primary).Debug("outbox record inserted")

	return nil
}
