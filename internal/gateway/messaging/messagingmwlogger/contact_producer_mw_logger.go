package messagingmwlogger

import (
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/model"

	"github.com/sirupsen/logrus"
)

type ContactProducerImpl struct {
	logger *logrus.Logger

	ProducerImpl[*model.ContactEvent]
	next messaging.ContactProducer
}

func NewContactProducer(logger *logrus.Logger, next messaging.ContactProducer) *ContactProducerImpl {
	return &ContactProducerImpl{
		ProducerImpl: ProducerImpl[*model.ContactEvent]{
			logger: logger,
			next:   next,
		},
		logger: logger,
		next:   next,
	}
}
