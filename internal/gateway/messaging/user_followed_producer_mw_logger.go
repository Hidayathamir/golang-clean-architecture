package messaging

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ UserFollowedProducer = &UserFollowedProducerMwLogger{}

type UserFollowedProducerMwLogger struct {
	Next UserFollowedProducer
}

func NewUserFollowedProducerMwLogger(next UserFollowedProducer) *UserFollowedProducerMwLogger {
	return &UserFollowedProducerMwLogger{
		Next: next,
	}
}

func (p *UserFollowedProducerMwLogger) Send(ctx context.Context, event *model.UserFollowedEvent) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := p.Next.Send(ctx, event)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"event": event,
	}
	x.LogMw(ctx, fields, err)

	return err
}
