package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/repository"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/google/uuid"
	"github.com/twmb/franz-go/pkg/kgo"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockOutboxProducer.go -pkg=mock . OutboxProducer

type OutboxProducer interface {
	ProducePending(ctx context.Context) error
}

var _ OutboxProducer = &OutboxProducerImpl{}

type OutboxProducerImpl struct {
	Cfg              *config.Config
	DB               *gorm.DB
	Client           *kgo.Client
	OutboxRepository repository.OutboxRepository
}

func NewOutboxProducer(
	cfg *config.Config,
	db *gorm.DB,
	client *kgo.Client,
	outboxRepository repository.OutboxRepository,
) *OutboxProducerImpl {
	return &OutboxProducerImpl{
		Cfg:              cfg,
		DB:               db,
		Client:           client,
		OutboxRepository: outboxRepository,
	}
}

func (p *OutboxProducerImpl) ProducePending(ctx context.Context) error {
	if p.Client == nil {
		logkit.Logger.WithContext(ctx).Warn("Kafka producer is disabled, skipping outbox produce")
		return nil
	}

	return p.DB.Transaction(func(tx *gorm.DB) error {
		outboxes := entity.OutboxList{}
		err := p.OutboxRepository.FindPending(ctx, tx, &outboxes, p.Cfg.GetOutboxBatchSize())
		if err != nil {
			return errkit.AddFuncName(err, "messaging.(*OutboxProducerImpl).ProducePending")
		}

		if len(outboxes) == 0 {
			return nil
		}

		var producedIDs []int64
		for _, outbox := range outboxes {
			recordCtx := telemetry.ExtractTraceContext(ctx, outbox.TraceContext)
			idempotencyKey := uuid.New().String()
			record := &kgo.Record{
				Topic: outbox.Topic,
				Value: outbox.Payload,
				Headers: []kgo.RecordHeader{
					{Key: "x-idempotency-key", Value: []byte(idempotencyKey)},
				},
			}

			result := p.Client.ProduceSync(recordCtx, record)
			err = result.FirstErr()
			if err != nil {
				logkit.Logger.WithContext(ctx).WithError(err).
					WithField("outbox_id", outbox.ID).
					Error("failed to produce outbox record to Kafka")
				continue
			}

			producedIDs = append(producedIDs, outbox.ID)
		}

		if len(producedIDs) > 0 {
			err = p.OutboxRepository.MarkProduced(ctx, tx, producedIDs)
			if err != nil {
				return errkit.AddFuncName(err, "messaging.(*OutboxProducerImpl).ProducePending")
			}
		}

		logkit.Logger.WithContext(ctx).
			WithField("produced", len(producedIDs)).
			WithField("total", len(outboxes)).
			Debug("outbox produce cycle complete")

		return nil
	})
}
