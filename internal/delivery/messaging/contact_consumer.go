package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/contact"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type ContactConsumer struct {
	Usecase contact.ContactUsecase
}

func NewContactConsumer(usecase contact.ContactUsecase) *ContactConsumer {
	return &ContactConsumer{
		Usecase: usecase,
	}
}

func (c ContactConsumer) Consume(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	ContactEvent := new(model.ContactEvent)
	if err := json.Unmarshal(message.Value, ContactEvent); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error("error unmarshalling Contact event")
		return errkit.AddFuncName(err)
	}

	// TODO process event
	x.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"event":     ContactEvent,
		"partition": message.Partition,
	}).Info("Received topic contacts")
	return nil
}
