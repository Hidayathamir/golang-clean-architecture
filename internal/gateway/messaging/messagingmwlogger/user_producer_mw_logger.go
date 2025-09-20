package messagingmwlogger

import (
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/model"

	"github.com/sirupsen/logrus"
)

type UserProducerImpl struct {
	logger *logrus.Logger

	ProducerImpl[*model.UserEvent]
	next messaging.UserProducer
}

func NewUserProducer(logger *logrus.Logger, next messaging.UserProducer) *UserProducerImpl {
	return &UserProducerImpl{
		ProducerImpl: ProducerImpl[*model.UserEvent]{
			logger: logger,
			next:   next,
		},
		logger: logger,
		next:   next,
	}
}
