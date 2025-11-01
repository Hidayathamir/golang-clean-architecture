package contact

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
)

var _ ContactUsecase = &ContactUsecaseMwTelemetry{}

type ContactUsecaseMwTelemetry struct {
	Next ContactUsecase
}

func NewContactUsecaseMwTelemetry(next ContactUsecase) *ContactUsecaseMwTelemetry {
	return &ContactUsecaseMwTelemetry{
		Next: next,
	}
}

func (u *ContactUsecaseMwTelemetry) Create(ctx context.Context, req *model.CreateContactRequest) (*model.ContactResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Create(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *ContactUsecaseMwTelemetry) Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Update(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *ContactUsecaseMwTelemetry) Get(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Get(ctx, req)
	telemetry.RecordError(span, err)

	return res, err
}

func (u *ContactUsecaseMwTelemetry) Delete(ctx context.Context, req *model.DeleteContactRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Delete(ctx, req)
	telemetry.RecordError(span, err)

	return err
}

func (u *ContactUsecaseMwTelemetry) Search(ctx context.Context, req *model.SearchContactRequest) (model.ContactResponseList, int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	list, total, err := u.Next.Search(ctx, req)
	telemetry.RecordError(span, err)

	return list, total, err
}
