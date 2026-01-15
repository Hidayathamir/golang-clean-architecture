package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type UserConsumer struct {
	Usecase user.UserUsecase
}

func NewUserConsumer(usecase user.UserUsecase) *UserConsumer {
	return &UserConsumer{
		Usecase: usecase,
	}
}

func (c UserConsumer) Consume(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	UserFollowedEvent := new(model.UserFollowedEvent)
	if err := json.Unmarshal(message.Value, UserFollowedEvent); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error("error unmarshalling User event")
		return errkit.AddFuncName(err)
	}

	// TODO process event
	x.Logger.WithContext(ctx).WithFields(logrus.Fields{
		"event":     UserFollowedEvent,
		"partition": message.Partition,
	}).Info("Received topic users")
	return nil
}
