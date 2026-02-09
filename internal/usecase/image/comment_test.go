package image_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/image"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestImageUsecaseImpl_Comment_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	CommentRepository := &mock.CommentRepositoryMock{}
	ImageProducer := &mock.ImageProducerMock{}

	u := &image.ImageUsecaseImpl{
		DB:                gormDB,
		CommentRepository: CommentRepository,
		ImageProducer:     ImageProducer,
	}

	req := &dto.CommentImageRequest{
		ImageID: 100,
		Comment: "nice pic",
	}

	CommentRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Comment) error {
		return nil
	}

	ImageProducer.SendImageCommentedFunc = func(ctx context.Context, event *dto.ImageCommentedEvent) error {
		return nil
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1})

	err := u.Comment(ctx, req)

	require.Nil(t, err)
}

func TestImageUsecaseImpl_Comment_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &image.ImageUsecaseImpl{
		DB: gormDB,
	}

	req := &dto.CommentImageRequest{} // invalid

	err := u.Comment(context.Background(), req)

	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestImageUsecaseImpl_Comment_Fail_Create(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	CommentRepository := &mock.CommentRepositoryMock{}

	u := &image.ImageUsecaseImpl{
		DB:                gormDB,
		CommentRepository: CommentRepository,
	}

	req := &dto.CommentImageRequest{
		ImageID: 100,
		Comment: "nice pic",
	}

	CommentRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Comment) error {
		return assert.AnError
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1})

	err := u.Comment(ctx, req)

	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestImageUsecaseImpl_Comment_Fail_Send(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	CommentRepository := &mock.CommentRepositoryMock{}
	ImageProducer := &mock.ImageProducerMock{}

	u := &image.ImageUsecaseImpl{
		DB:                gormDB,
		CommentRepository: CommentRepository,
		ImageProducer:     ImageProducer,
	}

	req := &dto.CommentImageRequest{
		ImageID: 100,
		Comment: "nice pic",
	}

	CommentRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Comment) error {
		return nil
	}

	ImageProducer.SendImageCommentedFunc = func(ctx context.Context, event *dto.ImageCommentedEvent) error {
		return assert.AnError
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1})

	err := u.Comment(ctx, req)

	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}
