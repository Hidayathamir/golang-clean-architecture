package messagingmwlogger

import (
	"context"
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ messaging.UserProducer = &UserProducerImpl{}

type UserProducerImpl struct {
	logger *logrus.Logger

	next messaging.UserProducer
}

func NewUserProducer(logger *logrus.Logger, next messaging.UserProducer) *UserProducerImpl {
	return &UserProducerImpl{
		logger: logger,
		next:   next,
	}
}

func (p *UserProducerImpl) Send(ctx context.Context, event *model.UserEvent) error {
	err := p.next.Send(ctx, event)

	fields := logrus.Fields{
		"event": event,
	}
	helper.Log(ctx, fields, err)

	return err
}
