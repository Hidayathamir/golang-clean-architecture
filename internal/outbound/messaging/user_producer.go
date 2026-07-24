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

//go:generate moq -out=../../mock/MockProducerUser2.go -pkg=mock . UserProducer

type UserProducer interface {
	SendUserFollowed(ctx context.Context, db *gorm.DB, event *dto.UserFollowedEvent) error
}

var _ UserProducer = &UserProducerImpl{}

type UserProducerImpl struct {
	Cfg              *config.Config
	OutboxRepository repository.OutboxRepository
}

func NewUserProducer(cfg *config.Config, outboxRepository repository.OutboxRepository) *UserProducerImpl {
	return &UserProducerImpl{
		Cfg:              cfg,
		OutboxRepository: outboxRepository,
	}
}

func (p *UserProducerImpl) SendUserFollowed(ctx context.Context, db *gorm.DB, event *dto.UserFollowedEvent) error {
	err := p.send(ctx, db, topic.UserFollowed, event)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*UserProducerImpl).SendUserFollowed")
	}
	return nil
}

func (p *UserProducerImpl) send(ctx context.Context, db *gorm.DB, topicName topic.Topic, event any) error {
	if !p.Cfg.GetKafkaProducerEnabled() {
		logkit.Logger.WithContext(ctx).Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*UserProducerImpl).send")
	}

	outbox := entity.Outbox{
		Topic:        topicName.Primary,
		Payload:      value,
		TraceContext: telemetry.InjectTraceContext(ctx),
		Status:       entity.OutboxStatusPending,
	}

	err = p.OutboxRepository.Insert(ctx, db, &outbox)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*UserProducerImpl).send")
	}

	logkit.Logger.WithContext(ctx).WithField("topic", topicName.Primary).Debug("outbox record inserted")

	return nil
}
