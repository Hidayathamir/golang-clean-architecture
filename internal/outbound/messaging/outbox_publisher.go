package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/repository"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/twmb/franz-go/pkg/kgo"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockOutboxPublisher.go -pkg=mock . OutboxPublisher

type OutboxPublisher interface {
	PublishPending(ctx context.Context) error
}

var _ OutboxPublisher = &OutboxPublisherImpl{}

type OutboxPublisherImpl struct {
	Cfg              *config.Config
	DB               *gorm.DB
	Client           *kgo.Client
	OutboxRepository repository.OutboxRepository
}

func NewOutboxPublisher(
	cfg *config.Config,
	db *gorm.DB,
	client *kgo.Client,
	outboxRepository repository.OutboxRepository,
) *OutboxPublisherImpl {
	return &OutboxPublisherImpl{
		Cfg:              cfg,
		DB:               db,
		Client:           client,
		OutboxRepository: outboxRepository,
	}
}

func (p *OutboxPublisherImpl) PublishPending(ctx context.Context) error {
	if p.Client == nil {
		logkit.Logger.WithContext(ctx).Warn("Kafka producer is disabled, skipping outbox publish")
		return nil
	}

	return p.DB.Transaction(func(tx *gorm.DB) error {
		outboxes := entity.OutboxList{}
		err := p.OutboxRepository.FindPending(ctx, tx, &outboxes, p.Cfg.GetOutboxBatchSize())
		if err != nil {
			return errkit.AddFuncName(err, "messaging.(*OutboxPublisherImpl).PublishPending")
		}

		if len(outboxes) == 0 {
			return nil
		}

		var publishedIDs []int64
		for _, outbox := range outboxes {
			record := &kgo.Record{
				Topic: outbox.Topic,
				Value: outbox.Payload,
			}

			result := p.Client.ProduceSync(ctx, record)
			err = result.FirstErr()
			if err != nil {
				logkit.Logger.WithContext(ctx).WithError(err).
					WithField("outbox_id", outbox.ID).
					Error("failed to publish outbox record to Kafka")
				continue
			}

			publishedIDs = append(publishedIDs, outbox.ID)
		}

		if len(publishedIDs) > 0 {
			err = p.OutboxRepository.MarkPublished(ctx, tx, publishedIDs)
			if err != nil {
				return errkit.AddFuncName(err, "messaging.(*OutboxPublisherImpl).PublishPending")
			}
		}

		logkit.Logger.WithContext(ctx).
			WithField("published", len(publishedIDs)).
			WithField("total", len(outboxes)).
			Debug("outbox publish cycle complete")

		return nil
	})
}
