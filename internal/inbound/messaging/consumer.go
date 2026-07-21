package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ConsumerHandlerBatch func(ctx context.Context, records []*kgo.Record) error

type ConsumerHandlerSingle func(ctx context.Context, record *kgo.Record) error

func ConsumeEventBatch(ctx context.Context, cfg *config.Config, consumerGroup string, topic string, handler ConsumerHandlerBatch) {
	localLogger := logkit.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"consumerGroup": consumerGroup,
		"topic":         topic,
	})

	localLogger.Info("setup kafka client")

	client := provider.NewKafkaClientConsumerBatch(cfg, consumerGroup, topic)

	for {
		const maxPollRecords = 0 // returns all buffered records
		fetches := client.PollRecords(ctx, maxPollRecords)

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				localLogger.WithError(err.Err).Error("error client poll fetch")
			}
			continue
		}

		records := fetches.Records()
		if len(records) > 0 {
			err := handler(ctx, records)
			if err != nil {
				localLogger.WithError(err).Error("handler got error processing message")
			} else {
				err = client.CommitUncommittedOffsets(ctx)
				if err != nil {
					localLogger.WithError(err).Error("client error commit")
				}
			}
		}

		if ctx.Err() != nil {
			localLogger.WithError(ctx.Err()).Info("context cancelled, stopping consumer")
			break
		}
	}

	localLogger.Info("Start closing consumer")
	client.Close()
	localLogger.Info("Done closing consumer")
}

func ConsumeEventSingle(ctx context.Context, cfg *config.Config, consumerGroup string, topic string, handler ConsumerHandlerSingle) {
	localLogger := logkit.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"consumerGroup": consumerGroup,
		"topic":         topic,
	})

	localLogger.Info("setup kafka client")

	client := provider.NewKafkaClientConsumerSingle(cfg, consumerGroup, topic)

	for {
		const maxPollRecords = 1 // returns a maximum 1 record.
		fetches := client.PollRecords(ctx, maxPollRecords)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				localLogger.WithError(err.Err).Error("error client poll fetch")
			}
			continue
		}

		records := fetches.Records()
		if len(records) > 0 {
			// its safe to directly use records[0] since maxPollRecords = 1
			err := handler(ctx, records[0])
			if err != nil {
				localLogger.WithError(err).Error("handler got error processing message")
			} else {
				err = client.CommitUncommittedOffsets(ctx)
				if err != nil {
					localLogger.WithError(err).Error("client error commit")
				}
			}
		}

		if ctx.Err() != nil {
			localLogger.WithError(ctx.Err()).Info("context cancelled, stopping consumer")
			break
		}
	}

	localLogger.Info("Start closing consumer")
	client.Close()
	localLogger.Info("Done closing consumer")
}
