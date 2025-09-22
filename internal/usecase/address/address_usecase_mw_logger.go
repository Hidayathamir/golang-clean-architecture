package address

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ AddressUseCase = &AddressUseCaseMwLogger{}

type AddressUseCaseMwLogger struct {
	logger *logrus.Logger

	next AddressUseCase
}

func NewAddressUseCaseMwLogger(logger *logrus.Logger, next AddressUseCase) *AddressUseCaseMwLogger {
	return &AddressUseCaseMwLogger{
		logger: logger,
		next:   next,
	}
}

func (u *AddressUseCaseMwLogger) Create(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error) {
	res, err := u.next.Create(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUseCaseMwLogger) Delete(ctx context.Context, req *model.DeleteAddressRequest) error {
	err := u.next.Delete(ctx, req)

	fields := logrus.Fields{
		"req": req,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (u *AddressUseCaseMwLogger) Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error) {
	res, err := u.next.Get(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUseCaseMwLogger) List(ctx context.Context, req *model.ListAddressRequest) ([]model.AddressResponse, error) {
	res, err := u.next.List(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *AddressUseCaseMwLogger) Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	res, err := u.next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}
