package image_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/image"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestImageUsecaseImpl_Like_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	LikeRepository := &mock.LikeRepositoryMock{}
	ImageProducer := &mock.ImageProducerMock{}

	u := &image.ImageUsecaseImpl{
		DB:             gormDB,
		LikeRepository: LikeRepository,
		ImageProducer:  ImageProducer,
	}

	req := &model.LikeImageRequest{
		ImageID: 100,
	}

	LikeRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Like) error {
		return nil
	}

	ImageProducer.SendImageLikedFunc = func(ctx context.Context, event *model.ImageLikedEvent) error {
		return nil
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &model.UserAuth{ID: 1})

	err := u.Like(ctx, req)

	require.Nil(t, err)
}

func TestImageUsecaseImpl_Like_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &image.ImageUsecaseImpl{
		DB: gormDB,
	}

	req := &model.LikeImageRequest{}

	err := u.Like(context.Background(), req)

	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestImageUsecaseImpl_Like_Fail_Create(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	LikeRepository := &mock.LikeRepositoryMock{}

	u := &image.ImageUsecaseImpl{
		DB:             gormDB,
		LikeRepository: LikeRepository,
	}

	req := &model.LikeImageRequest{
		ImageID: 100,
	}

	LikeRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Like) error {
		return assert.AnError
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &model.UserAuth{ID: 1})

	err := u.Like(ctx, req)

	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestImageUsecaseImpl_Like_Fail_Send(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	LikeRepository := &mock.LikeRepositoryMock{}
	ImageProducer := &mock.ImageProducerMock{}

	u := &image.ImageUsecaseImpl{
		DB:             gormDB,
		LikeRepository: LikeRepository,
		ImageProducer:  ImageProducer,
	}

	req := &model.LikeImageRequest{
		ImageID: 100,
	}

	LikeRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Like) error {
		return nil
	}

	ImageProducer.SendImageLikedFunc = func(ctx context.Context, event *model.ImageLikedEvent) error {
		return assert.AnError
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &model.UserAuth{ID: 1})

	err := u.Like(ctx, req)

	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}