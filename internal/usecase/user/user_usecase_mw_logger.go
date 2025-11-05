package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/l"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/sirupsen/logrus"
)

var _ UserUsecase = &UserUsecaseMwLogger{}

type UserUsecaseMwLogger struct {
	Next UserUsecase
}

func NewUserUsecaseMwLogger(next UserUsecase) *UserUsecaseMwLogger {
	return &UserUsecaseMwLogger{
		Next: next,
	}
}

func (u *UserUsecaseMwLogger) Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Create(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	l.LogMw(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Current(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	l.LogMw(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Login(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	l.LogMw(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Logout(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	l.LogMw(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Update(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	l.LogMw(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Verify(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	l.LogMw(ctx, fields, err)

	return res, err
}
