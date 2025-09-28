package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/helper"
	"github.com/sirupsen/logrus"
)

var _ AddressProducer = &AddressProducerMwLogger{}

type AddressProducerMwLogger struct {
	Next AddressProducer
}

func NewAddressProducerMwLogger(next AddressProducer) *AddressProducerMwLogger {
	return &AddressProducerMwLogger{
		Next: next,
	}
}

func (p *AddressProducerMwLogger) Send(ctx context.Context, event *model.AddressEvent) error {
	err := p.Next.Send(ctx, event)

	fields := logrus.Fields{
		"event": event,
	}
	helper.Log(ctx, fields, err)

	return err
}
