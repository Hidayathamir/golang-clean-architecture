package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ UserProducer = &UserProducerMwLogger{}

type UserProducerMwLogger struct {
	Next UserProducer
}

func NewUserProducerMwLogger(next UserProducer) *UserProducerMwLogger {
	return &UserProducerMwLogger{
		Next: next,
	}
}

func (p *UserProducerMwLogger) SendUserFollowed(ctx context.Context, event *dto.UserFollowedEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.SendUserFollowed(ctx, event)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	x.LogMw(ctx, fields, err)

	return err
}
