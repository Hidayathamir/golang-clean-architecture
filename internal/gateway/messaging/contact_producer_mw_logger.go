package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
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
	logging.Log(ctx, fields, err)

	return err
}
