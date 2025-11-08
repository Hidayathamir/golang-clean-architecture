package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ TodoProducer = &TodoProducerMwLogger{}

type TodoProducerMwLogger struct {
	Next TodoProducer
}

func NewTodoProducerMwLogger(next TodoProducer) *TodoProducerMwLogger {
	return &TodoProducerMwLogger{
		Next: next,
	}
}

func (p *TodoProducerMwLogger) Send(ctx context.Context, event *model.TodoCompletedEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.Send(ctx, event)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	x.LogMw(ctx, fields, err)

	return err
}
