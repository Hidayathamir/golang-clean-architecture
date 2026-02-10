package cache

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ UserCache = &UserCacheMwLogger{}

type UserCacheMwLogger struct {
	Next UserCache
}

func NewUserCacheMwLogger(next UserCache) *UserCacheMwLogger {
	return &UserCacheMwLogger{
		Next: next,
	}
}

func (u *UserCacheMwLogger) Get(ctx context.Context, id int64) (*entity.User, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	user, err := u.Next.Get(ctx, id)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id":   id,
		"user": user,
	}
	x.LogMw(ctx, fields, err)

	return user, err
}

func (u *UserCacheMwLogger) Set(ctx context.Context, user *entity.User) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Set(ctx, user)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"user": user,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (u *UserCacheMwLogger) Delete(ctx context.Context, id int64) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Delete(ctx, id)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"id": id,
	}
	x.LogMw(ctx, fields, err)

	return err
}
