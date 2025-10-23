package address

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
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
	logging.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) Delete(ctx context.Context, req *model.DeleteAddressRequest) error {
	err := u.Next.Delete(ctx, req)

	fields := logrus.Fields{
		"req": req,
	}
	logging.Log(ctx, fields, err)

	return err
}

func (u *AddressUsecaseMwLogger) Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error) {
	res, err := u.Next.Get(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	logging.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) List(ctx context.Context, req *model.ListAddressRequest) ([]model.AddressResponse, error) {
	res, err := u.Next.List(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	logging.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	res, err := u.Next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	logging.Log(ctx, fields, err)

	return res, err
}
