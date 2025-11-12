package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type TodoCompletionConsumer struct{}

func NewTodoCompletionConsumer() *TodoCompletionConsumer {
	return &TodoCompletionConsumer{}
}

func (c *TodoCompletionConsumer) Consume(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	event := new(model.TodoCompletedEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error("error unmarshalling todo completion event")
		return errkit.AddFuncName(err)
	}

	x.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"event":     event,
		"partition": message.Partition,
		"offset":    message.Offset,
	}).Info("Received todo completion event")

	return nil
}
