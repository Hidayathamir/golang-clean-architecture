package email

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ Sender = &SenderMwLogger{}

type SenderMwLogger struct {
	Next Sender
}

func NewSenderMwLogger(next Sender) *SenderMwLogger {
	return &SenderMwLogger{
		Next: next,
	}
}

func (s *SenderMwLogger) SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.EmailRetry(ctx, func() error {
		return s.Next.SendVerificationEmail(ctx, req)
	})
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"to": req.To,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}
