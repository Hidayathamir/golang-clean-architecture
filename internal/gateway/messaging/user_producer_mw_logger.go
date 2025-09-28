package messaging

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

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
	err := p.Next.Send(ctx, event)

	fields := logrus.Fields{
		"event": event,
	}
	helper.Log(ctx, fields, err)

	return err
}
