package repository

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ OutboxRepository = &OutboxRepositoryMwLogger{}

type OutboxRepositoryMwLogger struct {
	Next OutboxRepository
}

func NewOutboxRepositoryMwLogger(next OutboxRepository) *OutboxRepositoryMwLogger {
	return &OutboxRepositoryMwLogger{
		Next: next,
	}
}

func (r *OutboxRepositoryMwLogger) Insert(ctx context.Context, db *gorm.DB, outbox *entity.Outbox) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.Insert(ctx, db, outbox)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"topic":  outbox.Topic,
		"status": outbox.Status,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *OutboxRepositoryMwLogger) FindPending(ctx context.Context, db *gorm.DB, outboxes *entity.OutboxList, limit int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindPending(ctx, db, outboxes, limit)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"limit": limit,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (r *OutboxRepositoryMwLogger) MarkPublished(ctx context.Context, db *gorm.DB, ids []int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.MarkPublished(ctx, db, ids)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"ids": ids,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}
