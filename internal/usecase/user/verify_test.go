package user_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func newVerifyUsecase(t *testing.T) (*user.UserUsecaseImpl, *mock.UserRepositoryMock) {
	t.Helper()

	gormDB, _ := newFakeDB(t)
	repo := &mock.UserRepositoryMock{}

	cfg := config.NewConfig()
	cfg.SetAuthJWTSecret("test-secret")
	cfg.SetAuthJWTIssuer("test-issuer")
	cfg.SetAuthJWTExpireSeconds(60)

	u := &user.UserUsecaseImpl{
		Config:         cfg,
		DB:             gormDB,
		UserRepository: repo,
		UserCache:      newUserCacheMock(t),
	}

	return u, repo
}

func newSignedToken(t *testing.T, cfg *config.Config, userID int64) string {
	t.Helper()

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		Issuer:    cfg.GetAuthJWTIssuer(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.GetAuthJWTSecret()))
	require.NoError(t, err)
	return tokenString
}

func TestUserUsecaseImpl_Verify_Success(t *testing.T) {
	u, repo := newVerifyUsecase(t)

	tokenString := newSignedToken(t, u.Config, 1)

	req := &dto.VerifyUserRequest{
		Token: tokenString,
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id int64) error {
		entityMoqParam.ID = id
		return nil
	}

	res, err := u.Verify(context.Background(), *req)

	require.NoError(t, err)
	require.Equal(t, dto.UserAuth{ID: 1}, res)
}

func TestUserUsecaseImpl_Verify_ValidateStruct(t *testing.T) {
	u, repo := newVerifyUsecase(t)

	req := &dto.VerifyUserRequest{
		Token: "",
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id int64) error {
		return nil
	}

	res, err := u.Verify(context.Background(), *req)

	require.Equal(t, dto.UserAuth{}, res)
	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestUserUsecaseImpl_Verify_ParseAccessToken(t *testing.T) {
	u, repo := newVerifyUsecase(t)

	req := &dto.VerifyUserRequest{
		Token: "invalid-token",
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id int64) error {
		return nil
	}

	res, err := u.Verify(context.Background(), *req)

	require.Equal(t, dto.UserAuth{}, res)
	require.NotNil(t, err)
}

func TestUserUsecaseImpl_Verify_FindByID(t *testing.T) {
	u, repo := newVerifyUsecase(t)

	tokenString := newSignedToken(t, u.Config, 1)

	req := &dto.VerifyUserRequest{
		Token: tokenString,
	}

	repo.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.User, id int64) error {
		return assert.AnError
	}

	res, err := u.Verify(context.Background(), *req)

	require.Equal(t, dto.UserAuth{}, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}
