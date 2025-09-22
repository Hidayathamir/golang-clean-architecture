package messaging

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ UserProducer = &UserProducerMwLogger{}

type UserProducerMwLogger struct {
	logger *logrus.Logger

	next UserProducer
}

func NewUserProducerMwLogger(logger *logrus.Logger, next UserProducer) *UserProducerMwLogger {
	return &UserProducerMwLogger{
		logger: logger,
		next:   next,
	}
}

func (p *UserProducerMwLogger) Send(ctx context.Context, event *model.UserEvent) error {
	err := p.next.Send(ctx, event)

	fields := logrus.Fields{
		"event": event,
	}
	helper.Log(ctx, fields, err)

	return err
}
