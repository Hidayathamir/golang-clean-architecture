package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/IBM/sarama"
)

type UserConsumer struct {
	Usecase user.UserUsecase
}

func NewUserConsumer(usecase user.UserUsecase) *UserConsumer {
	return &UserConsumer{
		Usecase: usecase,
	}
}

func (c UserConsumer) ConsumeUserFollowedEvent(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	event := new(model.UserFollowedEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
		return errkit.AddFuncName(err)
	}

	req := new(model.NotifyUserBeingFollowedRequest)
	converter.ModelUserFollowedEventToModelNotifyUserBeingFollowedRequest(ctx, event, req)

	if err := c.Usecase.NotifyUserBeingFollowed(ctx, req); err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
