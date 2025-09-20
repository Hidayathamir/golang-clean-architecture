package messaging

import (
	"golang-clean-architecture/internal/model"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

//go:generate moq -out=../../mock/UserProducer.go -pkg=mock . UserProducer

type UserProducer interface {
	Producer[*model.UserEvent]
}

type UserProducerImpl struct {
	ProducerImpl[*model.UserEvent]
}

func NewUserProducer(producer sarama.SyncProducer, log *logrus.Logger) *UserProducerImpl {
	return &UserProducerImpl{
		ProducerImpl: ProducerImpl[*model.UserEvent]{
			Producer: producer,
			Topic:    "users",
			Log:      log,
		},
	}
}
