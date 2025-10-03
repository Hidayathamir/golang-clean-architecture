package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type UserConsumer struct {
	Usecase user.UserUsecase
	Log     *logrus.Logger
}

func NewUserConsumer(usecase user.UserUsecase, log *logrus.Logger) *UserConsumer {
	return &UserConsumer{
		Usecase: usecase,
		Log:     log,
	}
}

func (c UserConsumer) Consume(message *sarama.ConsumerMessage) error {
	UserEvent := new(model.UserEvent)
	if err := json.Unmarshal(message.Value, UserEvent); err != nil {
		c.Log.WithError(err).Error("error unmarshalling User event")
		return errkit.AddFuncName(err)
	}

	// TODO process event
	c.Log.WithFields(logrus.Fields{
		"event":     UserEvent,
		"partition": message.Partition,
	}).Info("Received topic users")
	return nil
}
