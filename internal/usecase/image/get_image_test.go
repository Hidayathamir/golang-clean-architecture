package image_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/image"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestImageUsecaseImpl_GetImage_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ImageRepository := &mock.ImageRepositoryMock{}

	u := &image.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: ImageRepository,
	}

	req := &dto.GetImageRequest{
		ID: 100,
	}

	ImageRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Image, id int64) error {
		entityMoqParam.ID = 100
		entityMoqParam.UserID = 1
		entityMoqParam.URL = "url"
		return nil
	}

	res, err := u.GetImage(context.Background(), req)

	expected := &dto.ImageResponse{
		ID:        100,
		UserID:    1,
		URL:       "url",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	require.Equal(t, expected, res)
	require.Nil(t, err)
}

func TestImageUsecaseImpl_GetImage_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &image.ImageUsecaseImpl{
		DB: gormDB,
	}

	req := &dto.GetImageRequest{}

	res, err := u.GetImage(context.Background(), req)

	require.Nil(t, res)
	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestImageUsecaseImpl_GetImage_Fail_FindByID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ImageRepository := &mock.ImageRepositoryMock{}

	u := &image.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: ImageRepository,
	}

	req := &dto.GetImageRequest{
		ID: 100,
	}

	ImageRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Image, id int64) error {
		return assert.AnError
	}

	res, err := u.GetImage(context.Background(), req)

	require.Nil(t, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}
