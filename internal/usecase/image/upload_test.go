package image_test

import (
	"context"
	"testing"
	"time"

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

func TestImageUsecaseImpl_Upload_Success(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ImageRepository := &mock.ImageRepositoryMock{}
	ImageProducer := &mock.ImageProducerMock{}
	S3Client := &mock.S3ClientMock{}

	u := &image.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: ImageRepository,
		ImageProducer:   ImageProducer,
		S3Client:        S3Client,
	}

	// ------------------------------------------------------- //

	req := &dto.UploadImageRequest{
		File: newFileHeader(t, "image.png", []byte("image-data")),
	}

	S3Client.UploadImageFunc = func(ctx context.Context, req dto.S3UploadImageRequest) (string, error) {
		return "http://image-url.com/image.png", nil
	}

	ImageRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Image) error {
		entityMoqParam.ID = 100
		return nil
	}

	ImageProducer.SendImageUploadedFunc = func(ctx context.Context, event *dto.ImageUploadedEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1, Username: "user1"})

	res, err := u.Upload(ctx, *req)

	// ------------------------------------------------------- //

	expected := dto.ImageResponse{
		ID:        100,
		UserID:    1,
		URL:       "http://image-url.com/image.png",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	require.Equal(t, expected, res)
	require.Nil(t, err)
}

func TestImageUsecaseImpl_Upload_Fail_ValidateStruct(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	u := &image.ImageUsecaseImpl{
		DB: gormDB,
	}

	req := &dto.UploadImageRequest{} // invalid request

	res, err := u.Upload(context.Background(), *req)

	require.Equal(t, dto.ImageResponse{}, res)
	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestImageUsecaseImpl_Upload_Fail_S3Upload(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	S3Client := &mock.S3ClientMock{}

	u := &image.ImageUsecaseImpl{
		DB:       gormDB,
		S3Client: S3Client,
	}

	req := &dto.UploadImageRequest{
		File: newFileHeader(t, "image.png", []byte("image-data")),
	}

	S3Client.UploadImageFunc = func(ctx context.Context, req dto.S3UploadImageRequest) (string, error) {
		return "", assert.AnError
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1, Username: "user1"})

	res, err := u.Upload(ctx, *req)

	require.Equal(t, dto.ImageResponse{}, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestImageUsecaseImpl_Upload_Fail_Create(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ImageRepository := &mock.ImageRepositoryMock{}
	S3Client := &mock.S3ClientMock{}

	u := &image.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: ImageRepository,
		S3Client:        S3Client,
	}

	req := &dto.UploadImageRequest{
		File: newFileHeader(t, "image.png", []byte("image-data")),
	}

	S3Client.UploadImageFunc = func(ctx context.Context, req dto.S3UploadImageRequest) (string, error) {
		return "url", nil
	}

	ImageRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Image) error {
		return assert.AnError
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1, Username: "user1"})

	res, err := u.Upload(ctx, *req)

	require.Equal(t, dto.ImageResponse{}, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestImageUsecaseImpl_Upload_Fail_Send(t *testing.T) {
	gormDB, _ := newFakeDB(t)
	ImageRepository := &mock.ImageRepositoryMock{}
	ImageProducer := &mock.ImageProducerMock{}
	S3Client := &mock.S3ClientMock{}

	u := &image.ImageUsecaseImpl{
		DB:              gormDB,
		ImageRepository: ImageRepository,
		ImageProducer:   ImageProducer,
		S3Client:        S3Client,
	}

	req := &dto.UploadImageRequest{
		File: newFileHeader(t, "image.png", []byte("image-data")),
	}

	S3Client.UploadImageFunc = func(ctx context.Context, req dto.S3UploadImageRequest) (string, error) {
		return "url", nil
	}

	ImageRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Image) error {
		return nil
	}

	ImageProducer.SendImageUploadedFunc = func(ctx context.Context, event *dto.ImageUploadedEvent) error {
		return assert.AnError
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1, Username: "user1"})

	res, err := u.Upload(ctx, *req)

	require.Equal(t, dto.ImageResponse{}, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}
