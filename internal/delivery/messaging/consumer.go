package messaging

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
)

type ConsumerHandler func(message *sarama.ConsumerMessage) error

type ConsumerGroupHandler struct {
	Handler ConsumerHandler
}

var _ sarama.ConsumerGroupHandler = &ConsumerGroupHandler{}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			err := h.Handler(message)
			if err != nil {
				x.Logger.WithError(err).Error("Failed to process message")
			} else {
				session.MarkMessage(message, "")
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

func ConsumeTopic(ctx context.Context, consumerGroup sarama.ConsumerGroup, topic string, handler ConsumerHandler) {
	consumerHandler := otelsarama.WrapConsumerGroupHandler(&ConsumerGroupHandler{
		Handler: handler,
	})

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, consumerHandler); err != nil {
				x.Logger.WithContext(ctx).WithError(err).Error("Error from consumer")
			}

			if ctx.Err() != nil {
				x.Logger.Info("Context cancelled, stopping consumer")
				return
			}
		}
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			x.Logger.WithError(err).Error("Consumer group error")
		}
	}()

	<-ctx.Done()
	x.Logger.Infof("Start closing consumer group for topic: %s", topic)
	if err := consumerGroup.Close(); err != nil {
		x.Logger.WithError(err).Error("Error closing consumer group")
	}
	x.Logger.Infof("Done closing consumer group for topic: %s", topic)
}

type BatchConsumerHandler func(messages []*sarama.ConsumerMessage) error

type BatchConsumerGroupHandler struct {
	Handler      BatchConsumerHandler
	BatchSize    int
	BatchTimeout time.Duration
}

var _ sarama.ConsumerGroupHandler = &BatchConsumerGroupHandler{}

func (h *BatchConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *BatchConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *BatchConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	batch := make([]*sarama.ConsumerMessage, 0, h.BatchSize)
	ticker := time.NewTicker(h.BatchTimeout)
	defer ticker.Stop()

	processBatch := func() {
		if len(batch) == 0 {
			return
		}
		if err := h.Handler(batch); err != nil {
			x.Logger.WithError(err).Error("Failed to process batch")
		} else {
			for _, msg := range batch {
				session.MarkMessage(msg, "")
			}
		}
		batch = batch[:0]
	}

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				processBatch()
				return nil
			}
			if message == nil {
				continue
			}
			batch = append(batch, message)
			if len(batch) >= h.BatchSize {
				processBatch()
			}
		case <-ticker.C:
			processBatch()
		case <-session.Context().Done():
			return nil
		}
	}
}

func ConsumeTopicBatch(ctx context.Context, consumerGroup sarama.ConsumerGroup, topic string, handler BatchConsumerHandler, batchSize int, batchTimeout time.Duration) {
	consumerHandler := otelsarama.WrapConsumerGroupHandler(&BatchConsumerGroupHandler{
		Handler:      handler,
		BatchSize:    batchSize,
		BatchTimeout: batchTimeout,
	})

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{topic}, consumerHandler); err != nil {
				x.Logger.WithContext(ctx).WithError(err).Error("Error from consumer")
			}

			if ctx.Err() != nil {
				x.Logger.Info("Context cancelled, stopping consumer")
				return
			}
		}
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			x.Logger.WithError(err).Error("Consumer group error")
		}
	}()

	<-ctx.Done()
	x.Logger.Infof("Start closing consumer group for topic: %s", topic)
	if err := consumerGroup.Close(); err != nil {
		x.Logger.WithError(err).Error("Error closing consumer group")
	}
	x.Logger.Infof("Done closing consumer group for topic: %s", topic)
}
