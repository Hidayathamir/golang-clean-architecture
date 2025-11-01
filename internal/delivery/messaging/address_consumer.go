package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type AddressConsumer struct {
	Usecase address.AddressUsecase
	Log     *logrus.Logger
}

func NewAddressConsumer(usecase address.AddressUsecase, log *logrus.Logger) *AddressConsumer {
	return &AddressConsumer{
		Usecase: usecase,
		Log:     log,
	}
}

func (c AddressConsumer) Consume(message *sarama.ConsumerMessage) error {
	_, span := telemetry.StartConsumer(message)
	defer span.End()

	addressEvent := new(model.AddressEvent)
	if err := json.Unmarshal(message.Value, addressEvent); err != nil {
		c.Log.WithError(err).Error("error unmarshalling address event")
		return errkit.AddFuncName("messaging.AddressConsumer.Consume", err)
	}

	// TODO process event
	c.Log.Infof("Received topic addresses with event: %v from partition %d", addressEvent, message.Partition)
	c.Log.WithFields(logrus.Fields{
		"event":     addressEvent,
		"partition": message.Partition,
	}).Info("Received topic address")
	return nil
}
