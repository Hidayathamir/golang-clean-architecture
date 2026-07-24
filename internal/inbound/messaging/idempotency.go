package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/idempotencyusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
)

func idempotencyKey(record *kgo.Record) string {
	for _, h := range record.Headers {
		if h.Key == "x-idempotency-key" {
			return string(h.Value)
		}
	}
	return ""
}

func IdempotencyHandlerSingle(usecase idempotencyusecase.IdempotencyUsecase, handler ConsumerHandlerSingle) ConsumerHandlerSingle {
	return func(ctx context.Context, record *kgo.Record) error {
		key := idempotencyKey(record)
		if key == "" {
			logkit.Logger.WithContext(ctx).Warn("record missing x-idempotency-key, processing without idempotency check")
			return handler(ctx, record)
		}

		isNew, err := usecase.InsertIfNotExists(ctx, key, record.Topic, record.Partition, record.Offset)
		if err != nil {
			logkit.Logger.WithContext(ctx).WithError(err).Error("idempotency check failed, processing without idempotency")
			return handler(ctx, record)
		}

		if !isNew {
			logkit.Logger.WithContext(ctx).WithFields(logrus.Fields{
				"idempotency_key": key,
				"topic":           record.Topic,
				"partition":       record.Partition,
				"offset":          record.Offset,
			}).Info("duplicate message skipped")
			return nil
		}

		return handler(ctx, record)
	}
}

func IdempotencyHandlerBatch(usecase idempotencyusecase.IdempotencyUsecase, handler ConsumerHandlerBatch) ConsumerHandlerBatch {
	return func(ctx context.Context, records []*kgo.Record) error {
		filtered := make([]*kgo.Record, 0, len(records))
		for _, record := range records {
			key := idempotencyKey(record)
			if key == "" {
				logkit.Logger.WithContext(ctx).Warn("record missing x-idempotency-key, including in batch without idempotency check")
				filtered = append(filtered, record)
				continue
			}

			isNew, err := usecase.InsertIfNotExists(ctx, key, record.Topic, record.Partition, record.Offset)
			if err != nil {
				logkit.Logger.WithContext(ctx).WithError(err).Error("idempotency check failed, including record without idempotency")
				filtered = append(filtered, record)
				continue
			}

			if !isNew {
				logkit.Logger.WithContext(ctx).WithFields(logrus.Fields{
					"idempotency_key": key,
					"topic":           record.Topic,
					"partition":       record.Partition,
					"offset":          record.Offset,
				}).Info("duplicate message skipped in batch")
				continue
			}

			filtered = append(filtered, record)
		}

		if len(filtered) == 0 {
			return nil
		}

		return handler(ctx, filtered)
	}
}
