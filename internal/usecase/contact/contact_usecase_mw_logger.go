package contact

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

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
	res, err := u.Next.Create(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *ContactUsecaseMwLogger) Delete(ctx context.Context, req *model.DeleteContactRequest) error {
	err := u.Next.Delete(ctx, req)

	fields := logrus.Fields{
		"req": req,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (u *ContactUsecaseMwLogger) Get(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error) {
	res, err := u.Next.Get(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *ContactUsecaseMwLogger) Search(ctx context.Context, req *model.SearchContactRequest) ([]model.ContactResponse, int64, error) {
	res, total, err := u.Next.Search(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, total, err
}

func (u *ContactUsecaseMwLogger) Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error) {
	res, err := u.Next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}
