package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/notifusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/twmb/franz-go/pkg/kgo"
)

type NotifConsumer struct {
	Usecase notifusecase.NotifUsecase
}

func NewNotifConsumer(usecase notifusecase.NotifUsecase) *NotifConsumer {
	return &NotifConsumer{
		Usecase: usecase,
	}
}

func (c *NotifConsumer) Notify(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.NotifEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(errkit.WrapNonRetryable(err), "messaging.(*NotifConsumer).Notify")
	}

	req := dto.NotifyRequest{}
	converter.DtoNotifEventToDtoNotifyRequest(event, &req)

	err = c.Usecase.Notify(ctx, req)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*NotifConsumer).Notify")
	}

	return nil
}
