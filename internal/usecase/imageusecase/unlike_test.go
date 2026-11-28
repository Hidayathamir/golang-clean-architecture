package imageusecase_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/imageusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestImageUsecaseImpl_Unlike_Success(t *testing.T) {
	gormDB, mockDB := newFakeDB(t)
	likeRepository := &mock.LikeRepositoryMock{}
	imageProducer := &mock.ImageProducerMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:             gormDB,
		LikeRepository: likeRepository,
		ImageProducer:  imageProducer,
	}

	likeRepository.DeleteByUserIDAndImageIDFunc = func(ctx context.Context, db *gorm.DB, userID int64, imageID int64) (int64, error) {
		require.Equal(t, int64(1), userID)
		require.Equal(t, int64(100), imageID)
		return 1, nil
	}

	var sentEvent *dto.ImageLikedEvent
	imageProducer.SendImageLikedFunc = func(ctx context.Context, db *gorm.DB, event *dto.ImageLikedEvent) error {
		sentEvent = event
		return nil
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 1})

	mockDB.ExpectBegin()
	mockDB.ExpectCommit()

	err := u.Unlike(ctx, dto.UnlikeImageRequest{ImageID: 100})

	require.Nil(t, err)
	require.NotNil(t, sentEvent)
	require.Equal(t, int64(100), sentEvent.ImageID)
	require.Equal(t, -1, sentEvent.Count)
}

func TestImageUsecaseImpl_Unlike_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &imageusecase.ImageUsecaseImpl{DB: gormDB}

	err := u.Unlike(context.Background(), dto.UnlikeImageRequest{ImageID: 0})

	require.NotNil(t, err)
}

func TestImageUsecaseImpl_Unlike_Fail_NotLiked(t *testing.T) {
	gormDB, mockDB := newFakeDB(t)
	likeRepository := &mock.LikeRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:             gormDB,
		LikeRepository: likeRepository,
		ImageProducer:  &mock.ImageProducerMock{},
	}

	likeRepository.DeleteByUserIDAndImageIDFunc = func(ctx context.Context, db *gorm.DB, userID int64, imageID int64) (int64, error) {
		return 0, nil
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 1})

	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	err := u.Unlike(ctx, dto.UnlikeImageRequest{ImageID: 100})

	require.NotNil(t, err)
	require.Equal(t, 404, errkit.GetHTTPError(err).HTTPCode)
}

func TestImageUsecaseImpl_Unlike_Fail_DeleteRepository(t *testing.T) {
	gormDB, mockDB := newFakeDB(t)
	likeRepository := &mock.LikeRepositoryMock{}

	u := &imageusecase.ImageUsecaseImpl{
		DB:             gormDB,
		LikeRepository: likeRepository,
		ImageProducer:  &mock.ImageProducerMock{},
	}

	likeRepository.DeleteByUserIDAndImageIDFunc = func(ctx context.Context, db *gorm.DB, userID int64, imageID int64) (int64, error) {
		return 0, assert.AnError
	}

	ctx := ctxuserauth.Set(context.Background(), &dto.UserAuth{ID: 1})

	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	err := u.Unlike(ctx, dto.UnlikeImageRequest{ImageID: 100})

	require.ErrorIs(t, err, assert.AnError)
}
