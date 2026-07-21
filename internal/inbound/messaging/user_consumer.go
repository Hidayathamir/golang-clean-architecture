package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/userusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/twmb/franz-go/pkg/kgo"
)

type UserConsumer struct {
	Usecase userusecase.UserUsecase
}

func NewUserConsumer(usecase userusecase.UserUsecase) *UserConsumer {
	return &UserConsumer{
		Usecase: usecase,
	}
}

func (c *UserConsumer) NotifyUserBeingFollowed(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.UserFollowedEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*UserConsumer).NotifyUserBeingFollowed")
	}

	req := dto.NotifyUserBeingFollowedRequest{}
	converter.DtoUserFollowedEventToDtoNotifyUserBeingFollowedRequest(event, &req)

	err = c.Usecase.NotifyUserBeingFollowed(ctx, req)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*UserConsumer).NotifyUserBeingFollowed")
	}

	return nil
}

func (c *UserConsumer) BatchUpdateUserFollowStats(ctx context.Context, records []*kgo.Record) error {
	ctx, span := telemetry.StartConsumerBatch(ctx, records)
	defer span.End()

	req := dto.BatchUpdateUserFollowStatsRequest{}
	converter.KGoRecordListToDtoBatchUpdateUserFollowStatsRequest(ctx, records, &req)

	err := c.Usecase.BatchUpdateUserFollowStats(ctx, req)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*UserConsumer).BatchUpdateUserFollowStats")
	}

	return nil
}
