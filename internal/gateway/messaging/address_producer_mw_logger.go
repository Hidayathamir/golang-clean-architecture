package messaging

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ AddressProducer = &AddressProducerMwLogger{}

type AddressProducerMwLogger struct {
	logger *logrus.Logger

	next AddressProducer
}

func NewAddressProducerMwLogger(logger *logrus.Logger, next AddressProducer) *AddressProducerMwLogger {
	return &AddressProducerMwLogger{
		logger: logger,
		next:   next,
	}
}

func (p *AddressProducerMwLogger) Send(ctx context.Context, event *model.AddressEvent) error {
	err := p.next.Send(ctx, event)

	fields := logrus.Fields{
		"event": event,
	}
	helper.Log(ctx, fields, err)

	return err
}
