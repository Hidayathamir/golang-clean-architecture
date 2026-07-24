package idempotencyusecase

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ IdempotencyUsecase = &IdempotencyUsecaseMwLogger{}

type IdempotencyUsecaseMwLogger struct {
	Next IdempotencyUsecase
}

func NewIdempotencyUsecaseMwLogger(next IdempotencyUsecase) *IdempotencyUsecaseMwLogger {
	return &IdempotencyUsecaseMwLogger{
		Next: next,
	}
}

func (u *IdempotencyUsecaseMwLogger) InsertIfNotExists(ctx context.Context, key, topic string, partition int32, offset int64) (bool, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	isNew, err := u.Next.InsertIfNotExists(ctx, key, topic, partition, offset)
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

func (u *IdempotencyUsecaseMwLogger) DeleteOlderThan(ctx context.Context, age time.Duration) (int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	deleted, err := u.Next.DeleteOlderThan(ctx, age)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"age":     age.String(),
		"deleted": deleted,
	}
	logkit.LogMw(ctx, fields, err)

	return deleted, err
}
