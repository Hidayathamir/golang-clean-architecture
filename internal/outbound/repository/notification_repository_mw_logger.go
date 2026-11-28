package repository

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/retrykit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var _ NotificationRepository = &NotificationRepositoryMwLogger{}

type NotificationRepositoryMwLogger struct {
	Next NotificationRepository
}

func NewNotificationRepositoryMwLogger(next NotificationRepository) *NotificationRepositoryMwLogger {
	return &NotificationRepositoryMwLogger{
		Next: next,
	}
}

func (r *NotificationRepositoryMwLogger) Create(ctx context.Context, db *gorm.DB, notification *entity.Notification) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.Create(ctx, db, notification)
	})
	telemetry.RecordError(span, err)
	fields := logrus.Fields{
		"notification": notification,
	}
	logkit.LogMw(ctx, fields, err)
	return err
}

func (r *NotificationRepositoryMwLogger) FindByID(ctx context.Context, db *gorm.DB, notification *entity.Notification, id int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindByID(ctx, db, notification, id)
	})
	telemetry.RecordError(span, err)
	fields := logrus.Fields{
		"id": id,
	}
	logkit.LogMw(ctx, fields, err)
	return err
}

func (r *NotificationRepositoryMwLogger) FindByUserID(ctx context.Context, db *gorm.DB, notificationList *entity.NotificationList, userID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := retrykit.DBRetry(ctx, func() error {
		return r.Next.FindByUserID(ctx, db, notificationList, userID, cursorCreatedAt, cursorID, limit)
	})
	telemetry.RecordError(span, err)
	fields := logrus.Fields{
		"user_id": userID,
		"limit":   limit,
	}
	logkit.LogMw(ctx, fields, err)
	return err
}

func (r *NotificationRepositoryMwLogger) MarkReadByID(ctx context.Context, db *gorm.DB, userID int64, notificationID int64) (int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	count, err := retrykit.DBRetryWithData(ctx, func() (int64, error) {
		return r.Next.MarkReadByID(ctx, db, userID, notificationID)
	})
	telemetry.RecordError(span, err)
	fields := logrus.Fields{
		"user_id":         userID,
		"notification_id": notificationID,
	}
	logkit.LogMw(ctx, fields, err)
	return count, err
}

func (r *NotificationRepositoryMwLogger) MarkAllReadByUserID(ctx context.Context, db *gorm.DB, userID int64) (int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	count, err := retrykit.DBRetryWithData(ctx, func() (int64, error) {
		return r.Next.MarkAllReadByUserID(ctx, db, userID)
	})
	telemetry.RecordError(span, err)
	fields := logrus.Fields{
		"user_id": userID,
	}
	logkit.LogMw(ctx, fields, err)
	return count, err
}

func (r *NotificationRepositoryMwLogger) CountUnreadByUserID(ctx context.Context, db *gorm.DB, userID int64) (int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	count, err := retrykit.DBRetryWithData(ctx, func() (int64, error) {
		return r.Next.CountUnreadByUserID(ctx, db, userID)
	})
	telemetry.RecordError(span, err)
	fields := logrus.Fields{
		"user_id": userID,
	}
	logkit.LogMw(ctx, fields, err)
	return count, err
}
