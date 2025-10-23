package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
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
	err := p.Next.Send(ctx, event)
	logging.Log(ctx, logrus.Fields{"event": event}, err)
	return err
}
