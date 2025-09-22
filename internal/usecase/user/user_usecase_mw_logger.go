package user

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ UserUseCase = &UserUseCaseMwLogger{}

type UserUseCaseMwLogger struct {
	logger *logrus.Logger

	next UserUseCase
}

func NewUserUseCaseMwLogger(logger *logrus.Logger, next UserUseCase) *UserUseCaseMwLogger {
	return &UserUseCaseMwLogger{
		logger: logger,
		next:   next,
	}
}

func (u *UserUseCaseMwLogger) Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Create(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *UserUseCaseMwLogger) Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Current(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *UserUseCaseMwLogger) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Login(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *UserUseCaseMwLogger) Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error) {
	res, err := u.next.Logout(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *UserUseCaseMwLogger) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *UserUseCaseMwLogger) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error) {
	res, err := u.next.Verify(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}
