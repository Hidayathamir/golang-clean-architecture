package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/notif"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/twmb/franz-go/pkg/kgo"
)

type NotifConsumer struct {
	Usecase notif.NotifUsecase
}

func NewNotifConsumer(usecase notif.NotifUsecase) *NotifConsumer {
	return &NotifConsumer{
		Usecase: usecase,
	}
}

func (c *NotifConsumer) Notify(ctx context.Context, records []*kgo.Record) error {
	for _, record := range records {
		err := c.notify(ctx, record)
		if err != nil {
			x.Logger.WithContext(ctx).WithError(err).Error()
			continue
		}
	}
	return nil
}

func (c *NotifConsumer) notify(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.NotifEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := dto.NotifyRequest{}
	converter.DtoNotifEventToDtoNotifyRequest(event, &req)

	err = c.Usecase.Notify(ctx, req)
	if err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}
