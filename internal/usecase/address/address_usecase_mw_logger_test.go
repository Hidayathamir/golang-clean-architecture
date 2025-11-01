package address_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logging"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewAddressUsecaseMwLogger(t *testing.T) {
	u := address.NewAddressUsecaseMwLogger(&mock.AddressUsecaseMock{})
	assert.NotEmpty(t, u)
}

func TestAddressUsecaseMwLogger_Create(t *testing.T) {
	logging.SetLogger(logrus.New())
	Next := &mock.AddressUsecaseMock{}
	u := &address.AddressUsecaseMwLogger{
		Next: Next,
	}
	Next.CreateFunc = func(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error) {
		return &model.AddressResponse{ID: "id1"}, nil
	}
	res, err := u.Create(context.Background(), &model.CreateAddressRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}

func TestAddressUsecaseMwLogger_Delete(t *testing.T) {
	logging.SetLogger(logrus.New())
	Next := &mock.AddressUsecaseMock{}
	u := &address.AddressUsecaseMwLogger{
		Next: Next,
	}
	Next.DeleteFunc = func(ctx context.Context, req *model.DeleteAddressRequest) error {
		return nil
	}
	err := u.Delete(context.Background(), &model.DeleteAddressRequest{})
	assert.Nil(t, err)
}

func TestAddressUsecaseMwLogger_Get(t *testing.T) {
	logging.SetLogger(logrus.New())
	Next := &mock.AddressUsecaseMock{}
	u := &address.AddressUsecaseMwLogger{
		Next: Next,
	}
	Next.GetFunc = func(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error) {
		return &model.AddressResponse{ID: "id1"}, nil
	}
	res, err := u.Get(context.Background(), &model.GetAddressRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}

func TestAddressUsecaseMwLogger_List(t *testing.T) {
	logging.SetLogger(logrus.New())
	Next := &mock.AddressUsecaseMock{}
	u := &address.AddressUsecaseMwLogger{
		Next: Next,
	}
	Next.ListFunc = func(ctx context.Context, req *model.ListAddressRequest) (model.AddressResponseList, error) {
		return model.AddressResponseList{{}}, nil
	}
	res, err := u.List(context.Background(), &model.ListAddressRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}

func TestAddressUsecaseMwLogger_Update(t *testing.T) {
	logging.SetLogger(logrus.New())
	Next := &mock.AddressUsecaseMock{}
	u := &address.AddressUsecaseMwLogger{
		Next: Next,
	}
	Next.UpdateFunc = func(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error) {
		return &model.AddressResponse{ID: "id1"}, nil
	}
	res, err := u.Update(context.Background(), &model.UpdateAddressRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}
