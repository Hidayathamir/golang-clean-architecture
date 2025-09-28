package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type AddressConsumer struct {
	Log *logrus.Logger
}

func NewAddressConsumer(log *logrus.Logger) *AddressConsumer {
	return &AddressConsumer{
		Log: log,
	}
}

func (c AddressConsumer) Consume(message *sarama.ConsumerMessage) error {
	addressEvent := new(model.AddressEvent)
	if err := json.Unmarshal(message.Value, addressEvent); err != nil {
		c.Log.WithError(err).Error("error unmarshalling address event")
		return errkit.AddFuncName(err)
	}

	// TODO process event
	c.Log.Infof("Received topic addresses with event: %v from partition %d", addressEvent, message.Partition)
	c.Log.WithFields(logrus.Fields{
		"event":     addressEvent,
		"partition": message.Partition,
	}).Info("Received topic address")
	return nil
}
