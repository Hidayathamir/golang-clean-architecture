package messaging

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ ContactProducer = &ContactProducerMwLogger{}

type ContactProducerMwLogger struct {
	logger *logrus.Logger

	next ContactProducer
}

func NewContactProducerMwLogger(logger *logrus.Logger, next ContactProducer) *ContactProducerMwLogger {
	return &ContactProducerMwLogger{
		logger: logger,
		next:   next,
	}
}

func (p *ContactProducerMwLogger) Send(ctx context.Context, event *model.ContactEvent) error {
	err := p.next.Send(ctx, event)

	fields := logrus.Fields{
		"event": event,
	}
	helper.Log(ctx, fields, err)

	return err
}
