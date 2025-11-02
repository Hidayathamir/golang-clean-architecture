package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/contact"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type ContactConsumer struct {
	Usecase contact.ContactUsecase
	Log     *logrus.Logger
}

func NewContactConsumer(usecase contact.ContactUsecase, log *logrus.Logger) *ContactConsumer {
	return &ContactConsumer{
		Usecase: usecase,
		Log:     log,
	}
}

func (c ContactConsumer) Consume(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	ContactEvent := new(model.ContactEvent)
	if err := json.Unmarshal(message.Value, ContactEvent); err != nil {
		c.Log.WithContext(ctx).WithError(err).Error("error unmarshalling Contact event")
		return errkit.AddFuncName("messaging.ContactConsumer.Consume", err)
	}

	// TODO process event
	c.Log.WithContext(ctx).WithFields(logrus.Fields{
		"event":     ContactEvent,
		"partition": message.Partition,
	}).Info("Received topic contacts")
	return nil
}
