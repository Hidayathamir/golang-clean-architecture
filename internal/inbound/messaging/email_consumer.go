package messaging

import (
	"context"
	"encoding/json"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/email"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/twmb/franz-go/pkg/kgo"
)

type EmailConsumer struct {
	EmailSender email.Sender
}

func NewEmailConsumer(emailSender email.Sender) *EmailConsumer {
	return &EmailConsumer{
		EmailSender: emailSender,
	}
}

func (c *EmailConsumer) SendVerificationEmail(ctx context.Context, record *kgo.Record) error {
	ctx, span := telemetry.StartConsumer(ctx, record)
	defer span.End()

	event := dto.EmailVerificationEvent{}
	err := json.Unmarshal(record.Value, &event)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(errkit.WrapNonRetryable(err), "messaging.(*EmailConsumer).SendVerificationEmail")
	}

	req := dto.SendVerificationEmailRequest{}
	converter.DtoEmailVerificationEventToDtoSendVerificationEmailRequest(event, &req)

	err = c.EmailSender.SendVerificationEmail(ctx, req)
	if err != nil {
		logkit.Logger.WithContext(ctx).WithError(err).Error()
		return errkit.AddFuncName(err, "messaging.(*EmailConsumer).SendVerificationEmail")
	}

	return nil
}
