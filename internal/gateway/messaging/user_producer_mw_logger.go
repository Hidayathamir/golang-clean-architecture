package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ UserProducer = &UserProducerMwLogger{}

type UserProducerMwLogger struct {
	Next UserProducer
}

func NewUserProducerMwLogger(next UserProducer) *UserProducerMwLogger {
	return &UserProducerMwLogger{
		Next: next,
	}
}

func (p *UserProducerMwLogger) Send(ctx context.Context, event *model.UserEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()
	err := p.Next.Send(ctx, event)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	logging.Log(ctx, fields, err)

	return err
}
