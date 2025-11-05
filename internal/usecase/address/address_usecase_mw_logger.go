package address

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/l"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
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
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Create(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	l.LogMw(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) Delete(ctx context.Context, req *model.DeleteAddressRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Delete(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	l.LogMw(ctx, fields, err)

	return err
}

func (u *AddressUsecaseMwLogger) Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Get(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	l.LogMw(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) List(ctx context.Context, req *model.ListAddressRequest) (model.AddressResponseList, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.List(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	l.LogMw(ctx, fields, err)

	return res, err
}

func (u *AddressUsecaseMwLogger) Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Update(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	l.LogMw(ctx, fields, err)

	return res, err
}
