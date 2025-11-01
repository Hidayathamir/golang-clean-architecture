package address

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
)

var _ AddressUsecase = &AddressUsecaseMwTelemetry{}

type AddressUsecaseMwTelemetry struct {
	Next AddressUsecase
}

func NewAddressUsecaseMwTelemetry(next AddressUsecase) *AddressUsecaseMwTelemetry {
	return &AddressUsecaseMwTelemetry{
		Next: next,
	}
}

func (u *AddressUsecaseMwTelemetry) Create(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Create(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *AddressUsecaseMwTelemetry) Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Update(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *AddressUsecaseMwTelemetry) Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Get(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *AddressUsecaseMwTelemetry) Delete(ctx context.Context, req *model.DeleteAddressRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Delete(ctx, req)
	telemetry.RecordError(span, err)

	return err
}

func (u *AddressUsecaseMwTelemetry) List(ctx context.Context, req *model.ListAddressRequest) (model.AddressResponseList, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	list, err := u.Next.List(ctx, req)
	telemetry.RecordError(span, err)

	return list, err
}
