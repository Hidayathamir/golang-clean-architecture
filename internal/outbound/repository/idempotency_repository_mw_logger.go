package repository

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ IdempotencyRepository = &IdempotencyRepositoryMwLogger{}

type IdempotencyRepositoryMwLogger struct {
	Next IdempotencyRepository
}

func NewIdempotencyRepositoryMwLogger(next IdempotencyRepository) *IdempotencyRepositoryMwLogger {
	return &IdempotencyRepositoryMwLogger{
		Next: next,
	}
}

func (r *IdempotencyRepositoryMwLogger) InsertIfNotExists(ctx context.Context, db *gorm.DB, key, topic string, partition int32, offset int64) (bool, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	var isNew bool
	err := retrykit.DBRetry(ctx, func() error {
		var innerErr error
		isNew, innerErr = r.Next.InsertIfNotExists(ctx, db, key, topic, partition, offset)
		return innerErr
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"idempotency_key": key,
		"topic":           topic,
		"partition":       partition,
		"record_offset":   offset,
		"is_new":          isNew,
	}
	logkit.LogMw(ctx, fields, err)

	return isNew, err
}

func (r *IdempotencyRepositoryMwLogger) DeleteOlderThan(ctx context.Context, db *gorm.DB, age time.Duration) (int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	var deleted int64
	err := retrykit.DBRetry(ctx, func() error {
		var innerErr error
		deleted, innerErr = r.Next.DeleteOlderThan(ctx, db, age)
		return innerErr
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"age":     age.String(),
		"deleted": deleted,
	}
	logkit.LogMw(ctx, fields, err)

	return deleted, err
}
