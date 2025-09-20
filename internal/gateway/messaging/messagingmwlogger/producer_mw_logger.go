package messagingmwlogger

import (
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ messaging.Producer[model.Event] = &ProducerImpl[model.Event]{}

type ProducerImpl[T model.Event] struct {
	logger *logrus.Logger

	next messaging.Producer[T]
}

func (p *ProducerImpl[T]) GetTopic() *string {
	return p.next.GetTopic()
}

func (p *ProducerImpl[T]) Send(event T) error {
	err := p.next.Send(event)

	fields := logrus.Fields{
		"event": event,
	}
	helper.Log(p.logger, fields, err)

	return err
}
