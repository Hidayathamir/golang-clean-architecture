package notif

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
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

func (u *NotifUsecaseMwLogger) Notify(ctx context.Context, req *model.NotifyRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Notify(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}
