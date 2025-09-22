package contact

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ ContactUsecase = &ContactUsecaseMwLogger{}

type ContactUsecaseMwLogger struct {
	logger *logrus.Logger

	next ContactUsecase
}

func NewContactUsecaseMwLogger(logger *logrus.Logger, next ContactUsecase) *ContactUsecaseMwLogger {
	return &ContactUsecaseMwLogger{
		logger: logger,
		next:   next,
	}
}

func (u *ContactUsecaseMwLogger) Create(ctx context.Context, req *model.CreateContactRequest) (*model.ContactResponse, error) {
	res, err := u.next.Create(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *ContactUsecaseMwLogger) Delete(ctx context.Context, req *model.DeleteContactRequest) error {
	err := u.next.Delete(ctx, req)

	fields := logrus.Fields{
		"req": req,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (u *ContactUsecaseMwLogger) Get(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error) {
	res, err := u.next.Get(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *ContactUsecaseMwLogger) Search(ctx context.Context, req *model.SearchContactRequest) ([]model.ContactResponse, int64, error) {
	res, total, err := u.next.Search(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, total, err
}

func (u *ContactUsecaseMwLogger) Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error) {
	res, err := u.next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}
