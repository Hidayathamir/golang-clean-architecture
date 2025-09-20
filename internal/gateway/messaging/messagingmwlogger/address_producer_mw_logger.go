package messagingmwlogger

import (
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ messaging.AddressProducer = &AddressProducerImpl{}

type AddressProducerImpl struct {
	logger *logrus.Logger

	next messaging.AddressProducer
}

func NewAddressProducer(logger *logrus.Logger, next messaging.AddressProducer) *AddressProducerImpl {
	return &AddressProducerImpl{
		logger: logger,
		next:   next,
	}
}

func (p *AddressProducerImpl) Send(event *model.AddressEvent) error {
	err := p.next.Send(event)

	fields := logrus.Fields{
		"event": event,
	}
	helper.Log(p.logger, fields, err)

	return err
}
