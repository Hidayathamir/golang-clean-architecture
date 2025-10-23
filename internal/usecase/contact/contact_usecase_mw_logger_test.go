package contact_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/contact"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewContactUsecaseMwLogger(t *testing.T) {
	u := contact.NewContactUsecaseMwLogger(&mock.ContactUsecaseMock{})
	assert.NotEmpty(t, u)
}

func TestContactUsecaseMwLogger_Create(t *testing.T) {
	logging.SetLogger(logrus.New())
	Next := &mock.ContactUsecaseMock{}
	u := &contact.ContactUsecaseMwLogger{
		Next: Next,
	}
	Next.CreateFunc = func(ctx context.Context, req *model.CreateContactRequest) (*model.ContactResponse, error) {
		return &model.ContactResponse{ID: "id1"}, nil
	}
	res, err := u.Create(context.Background(), &model.CreateContactRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}

func TestContactUsecaseMwLogger_Delete(t *testing.T) {
	logging.SetLogger(logrus.New())
	Next := &mock.ContactUsecaseMock{}
	u := &contact.ContactUsecaseMwLogger{
		Next: Next,
	}
	Next.DeleteFunc = func(ctx context.Context, req *model.DeleteContactRequest) error {
		return nil
	}
	err := u.Delete(context.Background(), &model.DeleteContactRequest{})
	assert.Nil(t, err)
}

func TestContactUsecaseMwLogger_Get(t *testing.T) {
	logging.SetLogger(logrus.New())
	Next := &mock.ContactUsecaseMock{}
	u := &contact.ContactUsecaseMwLogger{
		Next: Next,
	}
	Next.GetFunc = func(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error) {
		return &model.ContactResponse{ID: "id1"}, nil
	}
	res, err := u.Get(context.Background(), &model.GetContactRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}

func TestContactUsecaseMwLogger_Search(t *testing.T) {
	logging.SetLogger(logrus.New())
	Next := &mock.ContactUsecaseMock{}
	u := &contact.ContactUsecaseMwLogger{
		Next: Next,
	}
	Next.SearchFunc = func(ctx context.Context, req *model.SearchContactRequest) ([]model.ContactResponse, int64, error) {
		return []model.ContactResponse{{}}, 6, nil
	}
	res, total, err := u.Search(context.Background(), &model.SearchContactRequest{})
	assert.NotEmpty(t, res)
	assert.NotEmpty(t, total)
	assert.Nil(t, err)
}

func TestContactUsecaseMwLogger_Update(t *testing.T) {
	logging.SetLogger(logrus.New())
	Next := &mock.ContactUsecaseMock{}
	u := &contact.ContactUsecaseMwLogger{
		Next: Next,
	}
	Next.UpdateFunc = func(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error) {
		return &model.ContactResponse{ID: "id1"}, nil
	}
	res, err := u.Update(context.Background(), &model.UpdateContactRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}
