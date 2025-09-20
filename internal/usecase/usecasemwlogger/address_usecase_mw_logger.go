package usecasemwlogger

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ usecase.AddressUseCase = &AddressUseCaseImpl{}

type AddressUseCaseImpl struct {
	logger *logrus.Logger

	next usecase.AddressUseCase
}

func NewAddressUseCase(logger *logrus.Logger, next usecase.AddressUseCase) *AddressUseCaseImpl {
	return &AddressUseCaseImpl{
		logger: logger,
		next:   next,
	}
}

func (u *AddressUseCaseImpl) Create(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error) {
	res, err := u.next.Create(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}

func (u *AddressUseCaseImpl) Delete(ctx context.Context, req *model.DeleteAddressRequest) error {
	err := u.next.Delete(ctx, req)

	fields := logrus.Fields{
		"req": req,
	}
	helper.Log(u.logger, fields, err)

	return err
}

func (u *AddressUseCaseImpl) Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error) {
	res, err := u.next.Get(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}

func (u *AddressUseCaseImpl) List(ctx context.Context, req *model.ListAddressRequest) ([]model.AddressResponse, error) {
	res, err := u.next.List(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}

func (u *AddressUseCaseImpl) Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	res, err := u.next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}
