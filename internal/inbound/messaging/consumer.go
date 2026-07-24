package messaging

import (
	"context"
	"strconv"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/provider"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/topic"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
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

func produceToRetry(ctx context.Context, producer *kgo.Client, _topic topic.Topic, record *kgo.Record, retryCount int) {
	record.Topic = _topic.Retry()

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
			WithField("retryTopic", _topic.Retry()).
			WithField("retryCount", retryCount).
			Error("produceToRetry: failed to produce")
	}
}

func produceToDLQ(ctx context.Context, producer *kgo.Client, _topic topic.Topic, record *kgo.Record) {
	record.Topic = _topic.DLQ()

	result := producer.ProduceSync(ctx, record)
	if err := result.FirstErr(); err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).
			WithField("dlqTopic", _topic.DLQ()).
			Error("produceToDLQ: failed to produce")
	}
}

func ConsumeEventBatch(ctx context.Context, cfg *config.Config, producer *kgo.Client, consumerGroup string, _topic topic.Topic, handler ConsumerHandlerBatch) {
	localLogger := logkit.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"consumerGroup": consumerGroup,
		"topic":         _topic.Primary,
	})

	localLogger.Info("setup kafka client")

	client := provider.NewKafkaClientConsumerBatch(cfg, consumerGroup, _topic.Primary)

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

func ConsumeEventSingle(ctx context.Context, cfg *config.Config, producer *kgo.Client, consumerGroup string, _topic topic.Topic, handler ConsumerHandlerSingle) {
	localLogger := logkit.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"consumerGroup": consumerGroup,
		"topic":         _topic.Primary,
	})

	localLogger.Info("setup kafka client")

	client := provider.NewKafkaClientConsumerSingle(cfg, consumerGroup, _topic.Primary)

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

				if errkit.IsNonRetryable(err) {
					produceToDLQ(ctx, producer, _topic, record)
				} else {
					produceToRetry(ctx, producer, _topic, record, parseRetryCount(record)+1)
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

func ConsumeEventRetry(ctx context.Context, cfg *config.Config, producer *kgo.Client, consumerGroup string, _topic topic.Topic, handler ConsumerHandlerSingle) {
	localLogger := logkit.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"consumerGroup": consumerGroup,
		"topic":         _topic.Retry(),
	})

	localLogger.Info("setup kafka client")

	client := provider.NewKafkaClientConsumerSingle(cfg, consumerGroup, _topic.Retry())

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

				if errkit.IsNonRetryable(err) {
					produceToDLQ(ctx, producer, _topic, record)
				} else {
					retryCount := parseRetryCount(record)
					if retryCount >= maxRetries {
						produceToDLQ(ctx, producer, _topic, record)
					} else {
						produceToRetry(ctx, producer, _topic, record, retryCount+1)
					}
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
