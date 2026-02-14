package user_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNewUserUsecaseMwLogger(t *testing.T) {
	u := user.NewUserUsecaseMwLogger(&mock.UserUsecaseMock{})
	require.NotEmpty(t, u)
}

func TestUserUsecaseMwLogger_Create(t *testing.T) {
	x.SetLogger(logrus.New())
	Next := &mock.UserUsecaseMock{}
	u := &user.UserUsecaseMwLogger{
		Next: Next,
	}
	Next.CreateFunc = func(ctx context.Context, req dto.RegisterUserRequest) (dto.UserResponse, error) {
		return dto.UserResponse{ID: 1, Username: "user1"}, nil
	}
	res, err := u.Create(context.Background(), dto.RegisterUserRequest{})
	require.NotEmpty(t, res)
	require.Nil(t, err)
}

func TestUserUsecaseMwLogger_Current(t *testing.T) {
	x.SetLogger(logrus.New())
	Next := &mock.UserUsecaseMock{}
	u := &user.UserUsecaseMwLogger{
		Next: Next,
	}
	Next.CurrentFunc = func(ctx context.Context, req dto.GetUserRequest) (dto.UserResponse, error) {
		return dto.UserResponse{ID: 1, Username: "user1"}, nil
	}
	res, err := u.Current(context.Background(), dto.GetUserRequest{})
	require.NotEmpty(t, res)
	require.Nil(t, err)
}

func TestUserUsecaseMwLogger_Login(t *testing.T) {
	x.SetLogger(logrus.New())
	Next := &mock.UserUsecaseMock{}
	u := &user.UserUsecaseMwLogger{
		Next: Next,
	}
	Next.LoginFunc = func(ctx context.Context, req dto.LoginUserRequest) (dto.UserLoginResponse, error) {
		return dto.UserLoginResponse{ID: 1, Username: "user1"}, nil
	}
	res, err := u.Login(context.Background(), dto.LoginUserRequest{})
	require.NotEmpty(t, res)
	require.Nil(t, err)
}

func TestUserUsecaseMwLogger_Update(t *testing.T) {
	x.SetLogger(logrus.New())
	Next := &mock.UserUsecaseMock{}
	u := &user.UserUsecaseMwLogger{
		Next: Next,
	}
	Next.UpdateFunc = func(ctx context.Context, req dto.UpdateUserRequest) (dto.UserResponse, error) {
		return dto.UserResponse{ID: 1, Username: "user1"}, nil
	}
	res, err := u.Update(context.Background(), dto.UpdateUserRequest{})
	require.NotEmpty(t, res)
	require.Nil(t, err)
}

func TestUserUsecaseMwLogger_Verify(t *testing.T) {
	x.SetLogger(logrus.New())
	Next := &mock.UserUsecaseMock{}
	u := &user.UserUsecaseMwLogger{
		Next: Next,
	}
	Next.VerifyFunc = func(ctx context.Context, req dto.VerifyUserRequest) (dto.UserAuth, error) {
		return dto.UserAuth{ID: 1, Username: "user1"}, nil
	}
	res, err := u.Verify(context.Background(), dto.VerifyUserRequest{})
	require.NotEmpty(t, res)
	require.Nil(t, err)
}
