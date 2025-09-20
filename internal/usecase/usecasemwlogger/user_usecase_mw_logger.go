package usecasemwlogger

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ usecase.UserUseCase = &UserUseCaseImpl{}

type UserUseCaseImpl struct {
	logger *logrus.Logger

	next usecase.UserUseCase
}

func NewUserUseCase(logger *logrus.Logger, next usecase.UserUseCase) *UserUseCaseImpl {
	return &UserUseCaseImpl{
		logger: logger,
		next:   next,
	}
}

func (u *UserUseCaseImpl) Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Create(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}

func (u *UserUseCaseImpl) Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Current(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}

func (u *UserUseCaseImpl) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Login(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}

func (u *UserUseCaseImpl) Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error) {
	res, err := u.next.Logout(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}

func (u *UserUseCaseImpl) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	res, err := u.next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}

func (u *UserUseCaseImpl) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error) {
	res, err := u.next.Verify(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}
