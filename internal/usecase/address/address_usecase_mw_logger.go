package address

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ AddressUsecase = &AddressUsecaseMwLogger{}

type AddressUsecaseMwLogger struct {
	logger *logrus.Logger

	next AddressUsecase
}

func NewAddressUsecaseMwLogger(logger *logrus.Logger, next AddressUsecase) *AddressUsecaseMwLogger {
	return &AddressUsecaseMwLogger{
		logger: logger,
		next:   next,
	}
}

func (u *AddressUsecaseMwLogger) Create(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error) {
	res, err := u.next.Create(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) Delete(ctx context.Context, req *model.DeleteAddressRequest) error {
	err := u.next.Delete(ctx, req)

	fields := logrus.Fields{
		"req": req,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (u *AddressUsecaseMwLogger) Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error) {
	res, err := u.next.Get(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) List(ctx context.Context, req *model.ListAddressRequest) ([]model.AddressResponse, error) {
	res, err := u.next.List(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	res, err := u.next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}
