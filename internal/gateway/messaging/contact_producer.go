package messaging

import (
	"golang-clean-architecture/internal/model"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

//go:generate moq -out=../../mock/ContactProducer.go -pkg=mock . ContactProducer

type ContactProducer interface {
	Producer[*model.ContactEvent]
}

type ContactProducerImpl struct {
	ProducerImpl[*model.ContactEvent]
}

func NewContactProducer(producer sarama.SyncProducer, log *logrus.Logger) *ContactProducerImpl {
	return &ContactProducerImpl{
		ProducerImpl: ProducerImpl[*model.ContactEvent]{
			Producer: producer,
			Topic:    "contacts",
			Log:      log,
		},
	}
}
