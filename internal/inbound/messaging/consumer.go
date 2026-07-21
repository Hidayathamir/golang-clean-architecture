package messaging

import (
	"context"
	"strconv"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
)

type ConsumerHandlerBatch func(ctx context.Context, records []*kgo.Record) error

type ConsumerHandlerSingle func(ctx context.Context, record *kgo.Record) error

func parseRetryCount(record *kgo.Record) int {
	result := 0
	for _, h := range record.Headers {
		if h.Key == "x-retry-count" {
			n, err := strconv.Atoi(string(h.Value))
			if err == nil {
				result = n
			}
		}
	}
	return result
}

func produceToRetry(ctx context.Context, producer *kgo.Client, originalTopic string, record *kgo.Record, retryCount int) {
	retryTopic := topic.PrimaryToRetry[originalTopic]

	record.Topic = retryTopic

	filtered := make([]kgo.RecordHeader, 0, len(record.Headers))
	for _, h := range record.Headers {
		if h.Key != "x-retry-count" {
			filtered = append(filtered, h)
		}
	}
	filtered = append(filtered, kgo.RecordHeader{
		Key:   "x-retry-count",
		Value: []byte(strconv.Itoa(retryCount)),
	})
	record.Headers = filtered

	result := producer.ProduceSync(ctx, record)
	if err := result.FirstErr(); err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).
			WithField("retryTopic", retryTopic).
			WithField("retryCount", retryCount).
			Error("produceToRetry: failed to produce")
	}
}

func produceToDLQ(ctx context.Context, producer *kgo.Client, originalTopic string, record *kgo.Record) {
	dlqTopic := topic.PrimaryToDLQ[originalTopic]

	record.Topic = dlqTopic

	result := producer.ProduceSync(ctx, record)
	if err := result.FirstErr(); err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).
			WithField("dlqTopic", dlqTopic).
			Error("produceToDLQ: failed to produce")
	}
}

func ConsumeEventBatch(ctx context.Context, cfg *config.Config, producer *kgo.Client, consumerGroup string, _topic string, handler ConsumerHandlerBatch) {
	localLogger := logkit.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"consumerGroup": consumerGroup,
		"topic":         _topic,
	})

	localLogger.Info("setup kafka client")

	client := provider.NewKafkaClientConsumerBatch(cfg, consumerGroup, _topic)

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
				// TODO: improve by knowing which individual record failed in the batch,
				//       so only failed records go to retry
				for _, record := range records {
					produceToRetry(ctx, producer, _topic, record, 1)
				}
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

func ConsumeEventSingle(ctx context.Context, cfg *config.Config, producer *kgo.Client, consumerGroup string, primaryTopic string, handler ConsumerHandlerSingle) {
	localLogger := logkit.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"consumerGroup": consumerGroup,
		"topic":         primaryTopic,
	})

	localLogger.Info("setup kafka client")

	client := provider.NewKafkaClientConsumerSingle(cfg, consumerGroup, primaryTopic)

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
			record := records[0]
			err := handler(ctx, record)
			if err != nil {
				localLogger.WithError(err).Error("handler got error processing message")
				produceToRetry(ctx, producer, primaryTopic, record, 1)
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

func ConsumeEventRetry(ctx context.Context, cfg *config.Config, producer *kgo.Client, consumerGroup string, primaryTopic string, retryTopic string, handler ConsumerHandlerSingle) {
	localLogger := logkit.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"consumerGroup": consumerGroup,
		"topic":         primaryTopic,
	})

	localLogger.Info("setup kafka client")

	client := provider.NewKafkaClientConsumerSingle(cfg, consumerGroup, retryTopic)

	maxRetries := cfg.GetKafkaConsumerMaxRetries()

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
			record := records[0]
			err := handler(ctx, record)
			if err != nil {
				localLogger.WithError(err).Error("handler got error processing message")

				retryCount := parseRetryCount(record)
				if retryCount >= maxRetries {
					produceToDLQ(ctx, producer, primaryTopic, record)
				} else {
					produceToRetry(ctx, producer, primaryTopic, record, retryCount+1)
				}
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
