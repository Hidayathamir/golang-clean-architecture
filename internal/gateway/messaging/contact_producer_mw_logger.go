package messaging

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ ContactProducer = &ContactProducerMwLogger{}

type ContactProducerMwLogger struct {
	Next ContactProducer
}

func NewContactProducerMwLogger(next ContactProducer) *ContactProducerMwLogger {
	return &ContactProducerMwLogger{
		Next: next,
	}
}

func (p *ContactProducerMwLogger) Send(ctx context.Context, event *model.ContactEvent) error {
	err := p.Next.Send(ctx, event)

	fields := logrus.Fields{
		"event": event,
	}
	helper.Log(ctx, fields, err)

	return err
}
