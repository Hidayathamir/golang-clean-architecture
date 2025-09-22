package user

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ UserUsecase = &UserUsecaseMwLogger{}

type UserUsecaseMwLogger struct {
	logger *logrus.Logger

	next UserUsecase
}

func NewUserUsecaseMwLogger(logger *logrus.Logger, next UserUsecase) *UserUsecaseMwLogger {
	return &UserUsecaseMwLogger{
		logger: logger,
		next:   next,
	}
}

func (u *UserUsecaseMwLogger) Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Create(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Current(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Login(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error) {
	res, err := u.next.Logout(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *UserUsecaseMwLogger) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error) {
	res, err := u.next.Verify(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}
