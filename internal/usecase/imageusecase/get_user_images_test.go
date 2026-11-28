package imageusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/imageusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestImageUsecaseImpl_GetUserImages_Success_FirstPage(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	imageRepository := &mock.ImageRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: imageRepository,
	}

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	imageRepository.FindByUserIDFunc = func(ctx context.Context, db *gorm.DB, imageList *entity.ImageList, userID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		require.Equal(t, int64(5), userID)
		require.Equal(t, int64(0), cursorID)
		require.Equal(t, dto.DefaultPageLimit+1, limit)
		*imageList = entity.ImageList{
			{ID: 30, UserID: 5, URL: "u30", CreatedAt: base.Add(2 * time.Hour)},
			{ID: 20, UserID: 5, URL: "u20", CreatedAt: base.Add(1 * time.Hour)},
		}
		return nil
	}

	res, err := u.GetUserImages(context.Background(), dto.UserImagesRequest{UserID: 5})

	require.Nil(t, err)
	require.Len(t, res.Data, 2)
	require.Empty(t, res.NextCursor)
	require.Equal(t, int64(30), res.Data[0].ID)
}

func TestImageUsecaseImpl_GetUserImages_Success_NextPage(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	imageRepository := &mock.ImageRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: imageRepository,
	}

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cursor := pagination.Encode(pagination.Cursor{CreatedAt: base, ID: 100})

	imageRepository.FindByUserIDFunc = func(ctx context.Context, db *gorm.DB, imageList *entity.ImageList, userID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		require.Equal(t, int64(100), cursorID)
		require.True(t, base.Equal(cursorCreatedAt))
		require.Equal(t, 3, limit)
		*imageList = entity.ImageList{
			{ID: 3, UserID: 5, CreatedAt: base.Add(-time.Hour)},
			{ID: 2, UserID: 5, CreatedAt: base.Add(-2 * time.Hour)},
			{ID: 1, UserID: 5, CreatedAt: base.Add(-3 * time.Hour)},
		}
		return nil
	}

	req := dto.UserImagesRequest{UserID: 5, Page: dto.PageRequest{Limit: 2, Cursor: cursor}}
	res, err := u.GetUserImages(context.Background(), req)

	require.Nil(t, err)
	require.Len(t, res.Data, 2)

	next, ok, err := pagination.Decode(res.NextCursor)
	require.Nil(t, err)
	require.True(t, ok)
	require.Equal(t, int64(2), next.ID)
}

func TestImageUsecaseImpl_GetUserImages_Fail_InvalidCursor(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &imageusecase.ImageUsecaseImpl{DB: gormDB}

	req := dto.UserImagesRequest{UserID: 5, Page: dto.PageRequest{Cursor: "not-a-valid-cursor"}}
	res, err := u.GetUserImages(context.Background(), req)

	require.Equal(t, dto.UserImagesResponse{}, res)
	require.NotNil(t, err)
}

func TestImageUsecaseImpl_GetUserImages_Fail_Repository(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	imageRepository := &mock.ImageRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: imageRepository,
	}

	imageRepository.FindByUserIDFunc = func(ctx context.Context, db *gorm.DB, imageList *entity.ImageList, userID int64, cursorCreatedAt time.Time, cursorID int64, limit int) error {
		return assert.AnError
	}

	res, err := u.GetUserImages(context.Background(), dto.UserImagesRequest{UserID: 5})

	require.Equal(t, dto.UserImagesResponse{}, res)
	require.ErrorIs(t, err, assert.AnError)
}
