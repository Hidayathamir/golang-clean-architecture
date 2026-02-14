package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ConsumerHandler func(ctx context.Context, records []*kgo.Record) error

func ConsumeEvent(ctx context.Context, cfg *config.Config, consumerGroup string, topic string, handler ConsumerHandler) {
	localLogger := x.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"consumerGroup": consumerGroup,
		"topic":         topic,
	})

	localLogger.Info("setup kafka client")

	client := provider.NewKafkaClientConsumer(cfg, consumerGroup, topic)

	for {
		fetches := client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				localLogger.WithError(err.Err).Error("error client poll fetch")
			}
		}

		records := fetches.Records()
		if len(records) > 0 {
			err := handler(ctx, records)
			if err != nil {
				localLogger.WithError(err).Error("handler got error processing message")
			}

			err = client.CommitUncommittedOffsets(ctx)
			if err != nil {
				localLogger.WithError(err).Error("client error commit")
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
