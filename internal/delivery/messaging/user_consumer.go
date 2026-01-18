package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
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

func (c *UserConsumer) ConsumeUserFollowedEventForNotification(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	event := new(model.UserFollowedEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := new(model.NotifyUserBeingFollowedRequest)
	converter.ModelUserFollowedEventToModelNotifyUserBeingFollowedRequest(ctx, event, req)

	if err := c.Usecase.NotifyUserBeingFollowed(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *UserConsumer) ConsumeUserFollowedEventForUpdateCount(messages []*sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartNew()
	defer span.End()

	req := new(model.BatchUpdateUserFollowStatsRequest)
	converter.SaramaConsumerMessageListToModelBatchUpdateUserFollowStatsRequest(ctx, messages, req)

	if err := c.Usecase.BatchUpdateUserFollowStats(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}
