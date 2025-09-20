package usecasemwlogger

import (
	"context"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/usecase"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ usecase.ContactUseCase = &ContactUseCaseImpl{}

type ContactUseCaseImpl struct {
	logger *logrus.Logger

	next usecase.ContactUseCase
}

func NewContactUseCase(logger *logrus.Logger, next usecase.ContactUseCase) *ContactUseCaseImpl {
	return &ContactUseCaseImpl{
		logger: logger,
		next:   next,
	}
}

func (u *ContactUseCaseImpl) Create(ctx context.Context, req *model.CreateContactRequest) (*model.ContactResponse, error) {
	res, err := u.next.Create(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}

func (u *ContactUseCaseImpl) Delete(ctx context.Context, req *model.DeleteContactRequest) error {
	err := u.next.Delete(ctx, req)

	fields := logrus.Fields{
		"req": req,
	}
	helper.Log(u.logger, fields, err)

	return err
}

func (u *ContactUseCaseImpl) Get(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error) {
	res, err := u.next.Get(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}

func (u *ContactUseCaseImpl) Search(ctx context.Context, req *model.SearchContactRequest) ([]model.ContactResponse, int64, error) {
	res, total, err := u.next.Search(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, total, err
}

func (u *ContactUseCaseImpl) Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error) {
	res, err := u.next.Update(ctx, req)

	fields := logrus.Fields{
		"req": req,
		"res": res,
	}
	helper.Log(u.logger, fields, err)

	return res, err
}
