package notifusecase

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ NotifUsecase = &NotifUsecaseMwLogger{}

type NotifUsecaseMwLogger struct {
	Next NotifUsecase
}

func NewNotifUsecaseMwLogger(next NotifUsecase) *NotifUsecaseMwLogger {
	return &NotifUsecaseMwLogger{
		Next: next,
	}
}

func (u *NotifUsecaseMwLogger) Notify(ctx context.Context, req dto.NotifyRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Notify(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (u *NotifUsecaseMwLogger) List(ctx context.Context, req dto.NotificationsRequest) (dto.NotificationsResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.List(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	logkit.LogMw(ctx, fields, err)

	return res, err
}

func (u *NotifUsecaseMwLogger) MarkRead(ctx context.Context, req dto.MarkNotificationReadRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.MarkRead(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (u *NotifUsecaseMwLogger) MarkAllRead(ctx context.Context) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.MarkAllRead(ctx)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{}
	logkit.LogMw(ctx, fields, err)

	return err
}

func (u *NotifUsecaseMwLogger) CountUnread(ctx context.Context) (dto.UnreadCountResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.CountUnread(ctx)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{}
	logkit.LogMw(ctx, fields, err)

	return res, err
}
