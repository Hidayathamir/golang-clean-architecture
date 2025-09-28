package address

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ AddressUsecase = &AddressUsecaseMwLogger{}

type AddressUsecaseMwLogger struct {
	Next AddressUsecase
}

func NewAddressUsecaseMwLogger(next AddressUsecase) *AddressUsecaseMwLogger {
	return &AddressUsecaseMwLogger{
		Next: next,
	}
}

func (u *AddressUsecaseMwLogger) Create(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error) {
	res, err := u.Next.Create(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) Delete(ctx context.Context, req *model.DeleteAddressRequest) error {
	err := u.Next.Delete(ctx, req)

	fields := logrus.Fields{
		"req": req,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (u *AddressUsecaseMwLogger) Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error) {
	res, err := u.Next.Get(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) List(ctx context.Context, req *model.ListAddressRequest) ([]model.AddressResponse, error) {
	res, err := u.Next.List(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	res, err := u.Next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}
