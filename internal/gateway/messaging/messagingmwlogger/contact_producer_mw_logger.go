package messagingmwlogger

import (
	"context"
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ messaging.ContactProducer = &ContactProducerImpl{}

type ContactProducerImpl struct {
	logger *logrus.Logger

	next messaging.ContactProducer
}

func NewContactProducer(logger *logrus.Logger, next messaging.ContactProducer) *ContactProducerImpl {
	return &ContactProducerImpl{
		logger: logger,
		next:   next,
	}
}

func (p *ContactProducerImpl) Send(ctx context.Context, event *model.ContactEvent) error {
	err := p.next.Send(ctx, event)

	fields := logrus.Fields{
		"event": event,
	}
	helper.Log(ctx, fields, err)

	return err
}
