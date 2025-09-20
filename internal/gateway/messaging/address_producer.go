package messaging

import (
	"golang-clean-architecture/internal/model"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

//go:generate moq -out=../../mock/AddressProducer.go -pkg=mock . AddressProducer

type AddressProducer interface {
	Producer[*model.AddressEvent]
}

type AddressProducerImpl struct {
	ProducerImpl[*model.AddressEvent]
}

func NewAddressProducer(producer sarama.SyncProducer, log *logrus.Logger) *AddressProducerImpl {
	return &AddressProducerImpl{
		ProducerImpl: ProducerImpl[*model.AddressEvent]{
			Producer: producer,
			Topic:    "addresses",
			Log:      log,
		},
	}
}
