package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
)

type ConsumerHandler func(message *sarama.ConsumerMessage) error

type ConsumerGroupHandler struct {
	Handler ConsumerHandler
}

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
				x.Logger.WithError(err).Error("Error from consumer")
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
	x.Logger.Infof("Closing consumer group for topic: %s", topic)
	if err := consumerGroup.Close(); err != nil {
		x.Logger.WithError(err).Error("Error closing consumer group")
	}
}
