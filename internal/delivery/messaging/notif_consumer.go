package messaging

import (
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/notif"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/IBM/sarama"
)

type NotifConsumer struct {
	Usecase notif.NotifUsecase
}

func NewNotifConsumer(usecase notif.NotifUsecase) *NotifConsumer {
	return &NotifConsumer{
		Usecase: usecase,
	}
}

func (c *NotifConsumer) ConsumeNotifEvent(message *sarama.ConsumerMessage) error {
	ctx, span := telemetry.StartConsumer(message)
	defer span.End()

	event := new(dto.NotifEvent)
	if err := json.Unmarshal(message.Value, event); err != nil {
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
