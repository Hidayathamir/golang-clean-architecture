package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
)

var _ UserUsecase = &UserUsecaseMwTelemetry{}

type UserUsecaseMwTelemetry struct {
	Next UserUsecase
}

func NewUserUsecaseMwTelemetry(next UserUsecase) *UserUsecaseMwTelemetry {
	return &UserUsecaseMwTelemetry{
		Next: next,
	}
}

func (u *UserUsecaseMwTelemetry) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	auth, err := u.Next.Verify(ctx, req)
	telemetry.RecordError(span, err)

	return auth, err
}

func (u *UserUsecaseMwTelemetry) Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Create(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *UserUsecaseMwTelemetry) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Login(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *UserUsecaseMwTelemetry) Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Current(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *UserUsecaseMwTelemetry) Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	ok, err := u.Next.Logout(ctx, req)
	telemetry.RecordError(span, err)

	return ok, err
}

func (u *UserUsecaseMwTelemetry) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Update(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}
