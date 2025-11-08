package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type AddressConsumer struct {
	Usecase address.AddressUsecase
}

func NewAddressConsumer(usecase address.AddressUsecase) *AddressConsumer {
	return &AddressConsumer{
		Usecase: usecase,
	}
}

func (c AddressConsumer) Consume(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	addressEvent := new(model.AddressEvent)
	if err := json.Unmarshal(message.Value, addressEvent); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error("error unmarshalling address event")
		return errkit.AddFuncName("messaging.AddressConsumer.Consume", err)
	}

	// TODO process event
	x.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"event":     addressEvent,
		"partition": message.Partition,
	}).Info("Received topic address")
	return nil
}
