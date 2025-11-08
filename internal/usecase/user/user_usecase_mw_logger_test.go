package user_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewUserUsecaseMwLogger(t *testing.T) {
	u := user.NewUserUsecaseMwLogger(&mock.UserUsecaseMock{})
	assert.NotEmpty(t, u)
}

func TestUserUsecaseMwLogger_Create(t *testing.T) {
	x.SetLogger(logrus.New())
	Next := &mock.UserUsecaseMock{}
	u := &user.UserUsecaseMwLogger{
		Next: Next,
	}
	Next.CreateFunc = func(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
		return &model.UserResponse{ID: "id1"}, nil
	}
	res, err := u.Create(context.Background(), &model.RegisterUserRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}

func TestUserUsecaseMwLogger_Current(t *testing.T) {
	x.SetLogger(logrus.New())
	Next := &mock.UserUsecaseMock{}
	u := &user.UserUsecaseMwLogger{
		Next: Next,
	}
	Next.CurrentFunc = func(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
		return &model.UserResponse{ID: "id1"}, nil
	}
	res, err := u.Current(context.Background(), &model.GetUserRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}

func TestUserUsecaseMwLogger_Login(t *testing.T) {
	x.SetLogger(logrus.New())
	Next := &mock.UserUsecaseMock{}
	u := &user.UserUsecaseMwLogger{
		Next: Next,
	}
	Next.LoginFunc = func(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
		return &model.UserResponse{ID: "id1"}, nil
	}
	res, err := u.Login(context.Background(), &model.LoginUserRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}

func TestUserUsecaseMwLogger_Logout(t *testing.T) {
	x.SetLogger(logrus.New())
	Next := &mock.UserUsecaseMock{}
	u := &user.UserUsecaseMwLogger{
		Next: Next,
	}
	Next.LogoutFunc = func(ctx context.Context, req *model.LogoutUserRequest) (bool, error) {
		return true, nil
	}
	res, err := u.Logout(context.Background(), &model.LogoutUserRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}

func TestUserUsecaseMwLogger_Update(t *testing.T) {
	x.SetLogger(logrus.New())
	Next := &mock.UserUsecaseMock{}
	u := &user.UserUsecaseMwLogger{
		Next: Next,
	}
	Next.UpdateFunc = func(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
		return &model.UserResponse{ID: "id1"}, nil
	}
	res, err := u.Update(context.Background(), &model.UpdateUserRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}

func TestUserUsecaseMwLogger_Verify(t *testing.T) {
	x.SetLogger(logrus.New())
	Next := &mock.UserUsecaseMock{}
	u := &user.UserUsecaseMwLogger{
		Next: Next,
	}
	Next.VerifyFunc = func(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error) {
		return &model.Auth{ID: "id1"}, nil
	}
	res, err := u.Verify(context.Background(), &model.VerifyUserRequest{})
	assert.NotEmpty(t, res)
	assert.Nil(t, err)
}
