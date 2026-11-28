package cache

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ NotificationPublisher = &NotificationPublisherMwLogger{}

type NotificationPublisherMwLogger struct {
	Next NotificationPublisher
}

func NewNotificationPublisherMwLogger(next NotificationPublisher) *NotificationPublisherMwLogger {
	return &NotificationPublisherMwLogger{
		Next: next,
	}
}

func (p *NotificationPublisherMwLogger) Publish(ctx context.Context, userID int64, notification dto.NotificationResponse) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.CacheRetry(ctx, func() error {
		return p.Next.Publish(ctx, userID, notification)
	})
	telemetry.RecordError(span, err)
	fields := logrus.Fields{
		"user_id":      userID,
		"notification": notification,
	}
	logkit.LogMw(ctx, fields, err)
	return err
}
