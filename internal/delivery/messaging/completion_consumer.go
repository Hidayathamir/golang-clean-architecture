package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type TodoCompletionConsumer struct {
	Log *logrus.Logger
}

func NewTodoCompletionConsumer(log *logrus.Logger) *TodoCompletionConsumer {
	return &TodoCompletionConsumer{
		Log: log,
	}
}

func (c *TodoCompletionConsumer) Consume(message *sarama.ConsumerMessage) error {
	_, span := telemetry.StartConsumer(message)
	defer span.End()

	event := new(model.TodoCompletedEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
		c.Log.WithError(err).Error("error unmarshalling todo completion event")
		return errkit.AddFuncName("messaging.(*TodoCompletionConsumer).Consume", err)
	}

	c.Log.WithFields(logrus.Fields{
		"event":     event,
		"partition": message.Partition,
		"offset":    message.Offset,
	}).Info("Received todo completion event")

	return nil
}
