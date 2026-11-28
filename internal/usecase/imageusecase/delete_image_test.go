package imageusecase_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/imageusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestImageUsecaseImpl_DeleteImage_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	imageRepository := &mock.ImageRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: imageRepository,
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})

	imageRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error {
		require.Equal(t, int64(5), id)
		*image = entity.Image{ID: 5, UserID: 7}
		return nil
	}
	var deletedID int64
	imageRepository.DeleteByIDFunc = func(ctx context.Context, db *gorm.DB, image *entity.Image) error {
		deletedID = image.ID
		return nil
	}

	err := u.DeleteImage(ctx, dto.DeleteImageRequest{ID: 5})

	require.Nil(t, err)
	require.Equal(t, int64(5), deletedID)
}

func TestImageUsecaseImpl_DeleteImage_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &imageusecase.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: &mock.ImageRepositoryMock{},
	}

	err := u.DeleteImage(context.Background(), dto.DeleteImageRequest{ID: 0})

	require.NotNil(t, err)
}

func TestImageUsecaseImpl_DeleteImage_Fail_NotOwner(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	imageRepository := &mock.ImageRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: imageRepository,
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})

	imageRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error {
		*image = entity.Image{ID: 5, UserID: 99}
		return nil
	}

	err := u.DeleteImage(ctx, dto.DeleteImageRequest{ID: 5})

	require.NotNil(t, err)
	require.Equal(t, 403, errkit.GetHTTPError(err).HTTPCode)
}

func TestImageUsecaseImpl_DeleteImage_Fail_FindByID(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	imageRepository := &mock.ImageRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: imageRepository,
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 7})

	imageRepository.FindByIDFunc = func(ctx context.Context, db *gorm.DB, image *entity.Image, id int64) error {
		return assert.AnError
	}

	err := u.DeleteImage(ctx, dto.DeleteImageRequest{ID: 5})

	require.ErrorIs(t, err, assert.AnError)
}
