package contact

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
)

var _ ContactUsecase = &ContactUsecaseMwLogger{}

type ContactUsecaseMwLogger struct {
	Next ContactUsecase
}

func NewContactUsecaseMwLogger(next ContactUsecase) *ContactUsecaseMwLogger {
	return &ContactUsecaseMwLogger{
		Next: next,
	}
}

func (u *ContactUsecaseMwLogger) Create(ctx context.Context, req *model.CreateContactRequest) (*model.ContactResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Create(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}

func (u *ContactUsecaseMwLogger) Delete(ctx context.Context, req *model.DeleteContactRequest) error {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	err := u.Next.Delete(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
	}
	x.LogMw(ctx, fields, err)

	return err
}

func (u *ContactUsecaseMwLogger) Get(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Get(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}

func (u *ContactUsecaseMwLogger) Search(ctx context.Context, req *model.SearchContactRequest) (model.ContactResponseList, int64, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, total, err := u.Next.Search(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err)

	return res, total, err
}

func (u *ContactUsecaseMwLogger) Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	res, err := u.Next.Update(ctx, req)
	telemetry.RecordError(span, err)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	x.LogMw(ctx, fields, err)

	return res, err
}
