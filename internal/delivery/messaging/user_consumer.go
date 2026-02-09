package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/twmb/franz-go/pkg/kgo"
)

type UserConsumer struct {
	Usecase user.UserUsecase
}

func NewUserConsumer(usecase user.UserUsecase) *UserConsumer {
	return &UserConsumer{
		Usecase: usecase,
	}
}

func (c *UserConsumer) NotifyUserBeingFollowed(ctx context.Context, records []*kgo.Record) error {
	for _, record := range records {
		if err := c.notifyUserBeingFollowed(ctx, record); err != nil {
			x.Logger.WithContext(ctx).WithError(err).Error()
			continue
		}
	}
	return nil
}

func (c *UserConsumer) notifyUserBeingFollowed(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := new(dto.UserFollowedEvent)
	if err := json.Unmarshal(record.Value, event); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := new(dto.NotifyUserBeingFollowedRequest)
	converter.DtoUserFollowedEventToDtoNotifyUserBeingFollowedRequest(ctx, event, req)

	if err := c.Usecase.NotifyUserBeingFollowed(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *UserConsumer) BatchUpdateUserFollowStats(ctx context.Context, records []*kgo.Record) error {
	ctx, span := telemetry.StartConsumerBatch(ctx, records)
	defer span.End()

	req := new(dto.BatchUpdateUserFollowStatsRequest)
	converter.KGoRecordListToDtoBatchUpdateUserFollowStatsRequest(ctx, records, req)

	if err := c.Usecase.BatchUpdateUserFollowStats(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}
