package user_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func newLoginUsecase(t *testing.T) (*user.UserUsecaseImpl, *mock.UserRepositoryMock, *mock.UserProducerMock, *mock.SlackClientMock) {
	t.Helper()

	gormDB, _ := newFakeDB(t)
	repo := &mock.UserRepositoryMock{}
	producer := &mock.UserProducerMock{}
	slack := &mock.SlackClientMock{}

	cfg := viper.New()
	cfg.Set(configkey.AuthJWTSecret, "test-secret")
	cfg.Set(configkey.AuthJWTIssuer, "test-issuer")
	cfg.Set(configkey.AuthJWTExpireSeconds, 60)

	u := &user.UserUsecaseImpl{
		Config: cfg,
		DB:     gormDB,
		 UserRepository: repo,
		UserProducer: producer,
		SlackClient:  slack,
	}

	return u, repo, producer, slack
}

func TestUserUsecaseImpl_Login_Success(t *testing.T) {
	u, repo, producer, slack := newLoginUsecase(t)

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		entityMoqParam.ID = id
		return nil
	}

	slack.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	producer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	res, err := u.Login(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.NotEmpty(t, res.Token)

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(res.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.Config.GetString(configkey.AuthJWTSecret)), nil
	})
	require.NoError(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, req.ID, claims.Subject)
	assert.Equal(t, u.Config.GetString(configkey.AuthJWTIssuer), claims.Issuer)
	require.NotNil(t, claims.ExpiresAt)
	assert.WithinDuration(t, time.Now().Add(time.Minute), claims.ExpiresAt.Time, time.Minute)
}

func TestUserUsecaseImpl_Login_Fail_ValidateStruct(t *testing.T) {
	u, repo, producer, slack := newLoginUsecase(t)

	req := &model.LoginUserRequest{
		ID:       "",
		Password: "password1",
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		entityMoqParam.ID = id
		return nil
	}

	slack.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	producer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	res, err := u.Login(context.Background(), req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Login_Fail_FindByID(t *testing.T) {
	u, repo, producer, slack := newLoginUsecase(t)

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return assert.AnError
	}

	slack.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	producer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	res, err := u.Login(context.Background(), req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Login_Fail_CompareHashAndPassword(t *testing.T) {
	u, repo, producer, slack := newLoginUsecase(t)

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password2"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		entityMoqParam.ID = id
		return nil
	}

	slack.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	producer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	res, err := u.Login(context.Background(), req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestUserUsecaseImpl_Login_Fail_IsConnected(t *testing.T) {
	u, repo, producer, slack := newLoginUsecase(t)

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		entityMoqParam.ID = id
		return nil
	}

	slack.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return false, assert.AnError
	}

	producer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	res, err := u.Login(context.Background(), req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Login_Fail_Send(t *testing.T) {
	u, repo, producer, slack := newLoginUsecase(t)

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		entityMoqParam.ID = id
		return nil
	}

	slack.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	producer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return assert.AnError
	}

	res, err := u.Login(context.Background(), req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}

func TestUserUsecaseImpl_Login_Fail_SignAccessToken(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	repo := &mock.UserRepositoryMock{}
	producer := &mock.UserProducerMock{}
	slack := &mock.SlackClientMock{}

	cfg := viper.New()
	cfg.Set(configkey.AuthJWTIssuer, "test-issuer")
	cfg.Set(configkey.AuthJWTExpireSeconds, 60)

	u := &user.UserUsecaseImpl{
		Config: cfg,
		DB:     gormDB,
		 UserRepository: repo,
		UserProducer: producer,
		SlackClient:  slack,
	}

	req := &model.LoginUserRequest{
		ID:       "id1",
		Password: "password1",
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		pw, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
		require.NoError(t, err)
		entityMoqParam.Password = string(pw)
		entityMoqParam.ID = id
		return nil
	}

	slack.IsConnectedFunc = func(ctx context.Context) (bool, error) {
		return true, nil
	}

	producer.SendFunc = func(ctx context.Context, event *model.UserEvent) error {
		return nil
	}

	res, err := u.Login(context.Background(), req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
}
