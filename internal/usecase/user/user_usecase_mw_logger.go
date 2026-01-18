package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
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

func (u *UserUsecaseMwLogger) Follow(ctx context.Context, req *model.FollowUserRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Follow(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
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
	x.LogMw(ctx, fields, err)

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
	x.LogMw(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserLoginResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Login(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err)

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
	x.LogMw(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.UserAuth, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Verify(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) NotifyUserBeingFollowed(ctx context.Context, req *model.NotifyUserBeingFollowedRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.NotifyUserBeingFollowed(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (u *UserUsecaseMwLogger) BatchUpdateUserFollowStats(ctx context.Context, req *model.BatchUpdateUserFollowStatsRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.BatchUpdateUserFollowStats(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}
