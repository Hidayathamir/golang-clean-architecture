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
	"gorm.io/gorm"
)

func newVerifyUsecase(t *testing.T) (*user.UserUsecaseImpl, *mock.UserRepositoryMock) {
	t.Helper()

	gormDB, _ := newFakeDB(t)
	repo := &mock.UserRepositoryMock{}

	cfg := viper.New()
	cfg.Set(configkey.AuthJWTSecret, "test-secret")
	cfg.Set(configkey.AuthJWTIssuer, "test-issuer")
	cfg.Set(configkey.AuthJWTExpireSeconds, 60)

	u := &user.UserUsecaseImpl{
		Config:         cfg,
		DB:             gormDB,
		Validate:       validator.New(),
		UserRepository: repo,
	}

	return u, repo
}

func newSignedToken(t *testing.T, cfg *viper.Viper, userID string) string {
	t.Helper()

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    cfg.GetString(configkey.AuthJWTIssuer),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.GetString(configkey.AuthJWTSecret)))
	require.NoError(t, err)
	return tokenString
}

func TestUserUsecaseImpl_Verify_Success(t *testing.T) {
	u, repo := newVerifyUsecase(t)

	tokenString := newSignedToken(t, u.Config, "id1")

	req := &model.VerifyUserRequest{
		Token: tokenString,
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		entityMoqParam.ID = id
		return nil
	}

	res, err := u.Verify(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, &model.Auth{ID: "id1"}, res)
}

func TestUserUsecaseImpl_Verify_ValidateStruct(t *testing.T) {
	u, repo := newVerifyUsecase(t)

	req := &model.VerifyUserRequest{
		Token: "",
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	res, err := u.Verify(context.Background(), req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	var verrs validator.ValidationErrors
	assert.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Verify_ParseAccessToken(t *testing.T) {
	u, repo := newVerifyUsecase(t)

	req := &model.VerifyUserRequest{
		Token: "invalid-token",
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return nil
	}

	res, err := u.Verify(context.Background(), req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestUserUsecaseImpl_Verify_FindByID(t *testing.T) {
	u, repo := newVerifyUsecase(t)

	tokenString := newSignedToken(t, u.Config, "id1")

	req := &model.VerifyUserRequest{
		Token: tokenString,
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id string) error {
		return assert.AnError
	}

	res, err := u.Verify(context.Background(), req)

	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, assert.AnError)
}
