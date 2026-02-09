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
		if err := c.notify(ctx, record); err != nil {
			x.Logger.WithContext(ctx).WithError(err).Error()
			continue
		}
	}
	return nil
}

func (c *NotifConsumer) notify(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := new(dto.NotifEvent)
	if err := json.Unmarshal(record.Value, event); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	req := new(dto.NotifyRequest)
	converter.DtoNotifEventToDtoNotifyRequest(ctx, event, req)

	if err := c.Usecase.Notify(ctx, req); err != nil {
		x.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err)
	}

	return nil
}
