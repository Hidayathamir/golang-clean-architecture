package messagingmwlogger

import (
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/model"

	"github.com/sirupsen/logrus"
)

type AddressProducerImpl struct {
	logger *logrus.Logger

	ProducerImpl[*model.AddressEvent]
	next messaging.AddressProducer
}

func NewAddressProducer(logger *logrus.Logger, next messaging.AddressProducer) *AddressProducerImpl {
	return &AddressProducerImpl{
		ProducerImpl: ProducerImpl[*model.AddressEvent]{
			logger: logger,
			next:   next,
		},
		logger: logger,
		next:   next,
	}
}
