package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func newLoginUsecase(t *testing.T) (*user.UserUsecaseImpl, *mock.UserRepositoryMock, *mock.UserProducerMock) {
	t.Helper()

	gormDB, _ := newFakeDB(t)
	repo := &mock.UserRepositoryMock{}
	producer := &mock.UserProducerMock{}

	cfg := config.NewConfig()
	cfg.Set(configkey.AuthJWTSecret, "test-secret")
	cfg.Set(configkey.AuthJWTIssuer, "test-issuer")
	cfg.Set(configkey.AuthJWTExpireSeconds, 60)

	u := &user.UserUsecaseImpl{
		Config:         cfg,
		DB:             gormDB,
		UserRepository: repo,
		UserProducer:   producer,
	}

	return u, repo, producer
}

func TestUserUsecaseImpl_Login_Success(t *testing.T) {
	u, repo, producer := newLoginUsecase(t)

	req := &dto.LoginUserRequest{
		Username: "user1",
		Password: "password1",
	}

	repo.FindByUsernameFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, username string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		entityMoqParam.ID = 123
		entityMoqParam.Username = username
		return nil
	}

	producer.SendUserFollowedFunc = func(ctx context.Context, event *dto.UserFollowedEvent) error {
		return nil
	}

	res, err := u.Login(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.Token)

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(res.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.Config.GetString(configkey.AuthJWTSecret)), nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)
	require.Equal(t, "123", claims.Subject)
	require.Equal(t, u.Config.GetString(configkey.AuthJWTIssuer), claims.Issuer)
	require.NotNil(t, claims.ExpiresAt)
	assert.WithinDuration(t, time.Now().Add(time.Minute), claims.ExpiresAt.Time, time.Minute)
}

func TestUserUsecaseImpl_Login_Fail_ValidateStruct(t *testing.T) {
	u, repo, producer := newLoginUsecase(t)

	req := &dto.LoginUserRequest{
		Username: "",
		Password: "password1",
	}

	repo.FindByUsernameFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, username string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		entityMoqParam.ID = 123
		entityMoqParam.Username = username
		return nil
	}

	producer.SendUserFollowedFunc = func(ctx context.Context, event *dto.UserFollowedEvent) error {
		return nil
	}

	res, err := u.Login(context.Background(), req)

	require.Nil(t, res)
	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Login_Fail_FindByUsername(t *testing.T) {
	u, repo, producer := newLoginUsecase(t)

	req := &dto.LoginUserRequest{
		Username: "user1",
		Password: "password1",
	}

	repo.FindByUsernameFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, username string) error {
		return assert.AnError
	}

	producer.SendUserFollowedFunc = func(ctx context.Context, event *dto.UserFollowedEvent) error {
		return nil
	}

	res, err := u.Login(context.Background(), req)

	require.Nil(t, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Login_Fail_CompareHashAndPassword(t *testing.T) {
	u, repo, producer := newLoginUsecase(t)

	req := &dto.LoginUserRequest{
		Username: "user1",
		Password: "password1",
	}

	repo.FindByUsernameFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, username string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password2"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		entityMoqParam.ID = 123
		entityMoqParam.Username = username
		return nil
	}

	producer.SendUserFollowedFunc = func(ctx context.Context, event *dto.UserFollowedEvent) error {
		return nil
	}

	res, err := u.Login(context.Background(), req)

	require.Nil(t, res)
	require.NotNil(t, err)
}

func TestUserUsecaseImpl_Login_Fail_SignAccessToken(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	repo := &mock.UserRepositoryMock{}
	producer := &mock.UserProducerMock{}

	cfg := config.NewConfig()
	cfg.Set(configkey.AuthJWTIssuer, "test-issuer")
	cfg.Set(configkey.AuthJWTExpireSeconds, 60)

	u := &user.UserUsecaseImpl{
		Config:         cfg,
		DB:             gormDB,
		UserRepository: repo,
		UserProducer:   producer,
	}

	req := &dto.LoginUserRequest{
		Username: "user1",
		Password: "password1",
	}

	repo.FindByUsernameFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, username string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		entityMoqParam.ID = 123
		entityMoqParam.Username = username
		return nil
	}

	producer.SendUserFollowedFunc = func(ctx context.Context, event *dto.UserFollowedEvent) error {
		return nil
	}

	res, err := u.Login(context.Background(), req)

	require.Nil(t, res)
	require.NotNil(t, err)
}
