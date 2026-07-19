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

//go:generate moq -out=../../mock/MockProducerImage2.go -pkg=mock . ImageProducer

type ImageProducer interface {
	SendImageUploaded(ctx context.Context, db *gorm.DB, event *dto.ImageUploadedEvent) error
	SendImageLiked(ctx context.Context, db *gorm.DB, event *dto.ImageLikedEvent) error
	SendImageCommented(ctx context.Context, db *gorm.DB, event *dto.ImageCommentedEvent) error
}

var _ ImageProducer = &ImageProducerImpl{}

type ImageProducerImpl struct {
	Cfg             *config.Config
	OutboxRepository repository.OutboxRepository
}

func NewImageProducer(cfg *config.Config, outboxRepository repository.OutboxRepository) *ImageProducerImpl {
	return &ImageProducerImpl{
		Cfg:              cfg,
		OutboxRepository: outboxRepository,
	}
}

func (p *ImageProducerImpl) SendImageUploaded(ctx context.Context, db *gorm.DB, event *dto.ImageUploadedEvent) error {
	err := p.send(ctx, db, topic.ImageUploaded, event)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*ImageProducerImpl).SendImageUploaded")
	}
	return nil
}

func (p *ImageProducerImpl) SendImageLiked(ctx context.Context, db *gorm.DB, event *dto.ImageLikedEvent) error {
	err := p.send(ctx, db, topic.ImageLiked, event)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*ImageProducerImpl).SendImageLiked")
	}
	return nil
}

func (p *ImageProducerImpl) SendImageCommented(ctx context.Context, db *gorm.DB, event *dto.ImageCommentedEvent) error {
	err := p.send(ctx, db, topic.ImageCommented, event)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*ImageProducerImpl).SendImageCommented")
	}
	return nil
}

func (p *ImageProducerImpl) send(ctx context.Context, db *gorm.DB, topicName string, event any) error {
	if !p.Cfg.GetKafkaProducerEnabled() {
		logkit.Logger.WithContext(ctx).Warn("Kafka producer is disabled")
		return nil
	}

	value, err := json.Marshal(event)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*ImageProducerImpl).send")
	}

	outbox := entity.Outbox{
		Topic:        topicName,
		Payload:      value,
		TraceContext: telemetry.InjectTraceContext(ctx),
		Status:       entity.OutboxStatusPending,
	}

	err = p.OutboxRepository.Insert(ctx, db, &outbox)
	if err != nil {
		return errkit.AddFuncName(err, "messaging.(*ImageProducerImpl).send")
	}

	logkit.Logger.WithContext(ctx).WithField("topic", topicName).Debug("outbox record inserted")

	return nil
}
