package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ EmailProducer = &EmailProducerMwLogger{}

type EmailProducerMwLogger struct {
	Next EmailProducer
}

func NewEmailProducerMwLogger(next EmailProducer) *EmailProducerMwLogger {
	return &EmailProducerMwLogger{
		Next: next,
	}
}

func (p *EmailProducerMwLogger) SendVerificationEmail(ctx context.Context, db *gorm.DB, event *dto.EmailVerificationEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.MessagingRetry(ctx, func() error {
		return p.Next.SendVerificationEmail(ctx, db, event)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}
