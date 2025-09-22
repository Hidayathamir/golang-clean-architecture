package contact

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ ContactUseCase = &ContactUseCaseMwLogger{}

type ContactUseCaseMwLogger struct {
	logger *logrus.Logger

	next ContactUseCase
}

func NewContactUseCaseMwLogger(logger *logrus.Logger, next ContactUseCase) *ContactUseCaseMwLogger {
	return &ContactUseCaseMwLogger{
		logger: logger,
		next:   next,
	}
}

func (u *ContactUseCaseMwLogger) Create(ctx context.Context, req *model.CreateContactRequest) (*model.ContactResponse, error) {
	res, err := u.next.Create(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *ContactUseCaseMwLogger) Delete(ctx context.Context, req *model.DeleteContactRequest) error {
	err := u.next.Delete(ctx, req)

	fields := logrus.Fields{
		"req": req,
	}
	helper.Log(ctx, fields, err)

	return err
}

func (u *ContactUseCaseMwLogger) Get(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error) {
	res, err := u.next.Get(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}

func (u *ContactUseCaseMwLogger) Search(ctx context.Context, req *model.SearchContactRequest) ([]model.ContactResponse, int64, error) {
	res, total, err := u.next.Search(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, total, err
}

func (u *ContactUseCaseMwLogger) Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error) {
	res, err := u.next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(ctx, fields, err)

	return res, err
}
