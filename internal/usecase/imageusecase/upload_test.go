package imageusecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/imageusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestImageUsecaseImpl_PresignUpload_Success(t *testing.T) {
	cfg := config.NewConfig()
	s3Client := &mock.S3ClientMock{}

	u := &imageusecase.ImageUsecaseImpl{
		Cfg:      cfg,
		S3Client: s3Client,
	}

	// ------------------------------------------------------- //

	req := &dto.PresignUploadImageRequest{
		Filename:    "image.png",
		ContentType: "image/png",
	}

	var got dto.S3PresignUploadRequest
	s3Client.PresignUploadImageFunc = func(ctx context.Context, req dto.S3PresignUploadRequest) (dto.S3PresignUploadResponse, error) {
		got = req
		return dto.S3PresignUploadResponse{UploadURL: "http://presigned.url", Key: req.Key}, nil
	}

	// ------------------------------------------------------- //

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1, Username: "user1"})

	res, err := u.PresignUpload(ctx, *req)

	// ------------------------------------------------------- //

	require.Nil(t, err)
	require.Equal(t, "image/png", got.ContentType)
	require.NotEmpty(t, got.Key)
	require.Equal(t, "http://presigned.url", res.UploadURL)
	require.Equal(t, got.Key, res.Key)
	require.Contains(t, res.URL, got.Key)
}

func TestImageUsecaseImpl_PresignUpload_Fail_ValidateStruct(t *testing.T) {
	cfg := config.NewConfig()
	u := &imageusecase.ImageUsecaseImpl{
		Cfg: cfg,
	}

	req := &dto.PresignUploadImageRequest{} // invalid request

	res, err := u.PresignUpload(context.Background(), *req)

	require.Equal(t, dto.PresignUploadImageResponse{}, res)
	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestImageUsecaseImpl_PresignUpload_Fail_Presign(t *testing.T) {
	cfg := config.NewConfig()
	s3Client := &mock.S3ClientMock{}

	u := &imageusecase.ImageUsecaseImpl{
		Cfg:      cfg,
		S3Client: s3Client,
	}

	req := &dto.PresignUploadImageRequest{
		Filename:    "image.png",
		ContentType: "image/png",
	}

	s3Client.PresignUploadImageFunc = func(ctx context.Context, req dto.S3PresignUploadRequest) (dto.S3PresignUploadResponse, error) {
		return dto.S3PresignUploadResponse{}, assert.AnError
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1, Username: "user1"})

	res, err := u.PresignUpload(ctx, *req)

	require.Equal(t, dto.PresignUploadImageResponse{}, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestImageUsecaseImpl_ConfirmUpload_Success(t *testing.T) {
	cfg := config.NewConfig()
	gormDB, mockDB := newFakeDB(t)
	ImageRepository := &mock.ImageRepositoryMock{}
	ImageProducer := &mock.ImageProducerMock{}
	S3Client := &mock.S3ClientMock{}

	u := &imageusecase.ImageUsecaseImpl{
		Cfg:             cfg,
		DB:              gormDB,
		ImageRepository: ImageRepository,
		ImageProducer:   ImageProducer,
		S3Client:        S3Client,
	}

	// ------------------------------------------------------- //

	req := &dto.ConfirmUploadImageRequest{
		Key:     "user1/123_image.png",
		Caption: "my caption",
	}

	S3Client.ObjectExistsFunc = func(ctx context.Context, req dto.S3ObjectExistsRequest) (dto.S3ObjectExistsResponse, error) {
		return dto.S3ObjectExistsResponse{Exists: true}, nil
	}

	ImageRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Image) error {
		entityMoqParam.ID = 100
		return nil
	}

	ImageProducer.SendImageUploadedFunc = func(ctx context.Context, db *gorm.DB, event *dto.ImageUploadedEvent) error {
		return nil
	}

	// ------------------------------------------------------- //

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1, Username: "user1"})

	mockDB.ExpectBegin()
	mockDB.ExpectCommit()

	res, err := u.ConfirmUpload(ctx, *req)

	// ------------------------------------------------------- //

	expected := dto.ImageResponse{
		ID:        100,
		UserID:    1,
		Caption:   "my caption",
		URL:       cfg.GetAWSBaseEndpoint() + "/image/" + req.Key,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	}

	require.Equal(t, expected, res)
	require.Nil(t, err)
}

func TestImageUsecaseImpl_ConfirmUpload_Fail_ValidateStruct(t *testing.T) {
	cfg := config.NewConfig()
	gormDB, _ := newFakeDB(t)
	u := &imageusecase.ImageUsecaseImpl{
		Cfg: cfg,
		DB:  gormDB,
	}

	req := &dto.ConfirmUploadImageRequest{} // invalid request

	res, err := u.ConfirmUpload(context.Background(), *req)

	require.Equal(t, dto.ImageResponse{}, res)
	require.NotNil(t, err)
	var verrs validator.ValidationErrors
	require.ErrorAs(t, err, &verrs)
}

func TestImageUsecaseImpl_ConfirmUpload_Fail_ObjectMissing(t *testing.T) {
	cfg := config.NewConfig()
	gormDB, _ := newFakeDB(t)
	S3Client := &mock.S3ClientMock{}

	u := &imageusecase.ImageUsecaseImpl{
		Cfg:      cfg,
		DB:       gormDB,
		S3Client: S3Client,
	}

	req := &dto.ConfirmUploadImageRequest{Key: "user1/123_image.png"}

	S3Client.ObjectExistsFunc = func(ctx context.Context, req dto.S3ObjectExistsRequest) (dto.S3ObjectExistsResponse, error) {
		return dto.S3ObjectExistsResponse{Exists: false}, nil
	}

	res, err := u.ConfirmUpload(context.Background(), *req)

	require.Equal(t, dto.ImageResponse{}, res)
	require.NotNil(t, err)
}

func TestImageUsecaseImpl_ConfirmUpload_Fail_ObjectExists(t *testing.T) {
	cfg := config.NewConfig()
	gormDB, _ := newFakeDB(t)
	S3Client := &mock.S3ClientMock{}

	u := &imageusecase.ImageUsecaseImpl{
		Cfg:      cfg,
		DB:       gormDB,
		S3Client: S3Client,
	}

	req := &dto.ConfirmUploadImageRequest{Key: "user1/123_image.png"}

	S3Client.ObjectExistsFunc = func(ctx context.Context, req dto.S3ObjectExistsRequest) (dto.S3ObjectExistsResponse, error) {
		return dto.S3ObjectExistsResponse{}, assert.AnError
	}

	res, err := u.ConfirmUpload(context.Background(), *req)

	require.Equal(t, dto.ImageResponse{}, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestImageUsecaseImpl_ConfirmUpload_Fail_Create(t *testing.T) {
	cfg := config.NewConfig()
	gormDB, mockDB := newFakeDB(t)
	ImageRepository := &mock.ImageRepositoryMock{}
	S3Client := &mock.S3ClientMock{}

	u := &imageusecase.ImageUsecaseImpl{
		Cfg:             cfg,
		DB:              gormDB,
		ImageRepository: ImageRepository,
		S3Client:        S3Client,
	}

	req := &dto.ConfirmUploadImageRequest{Key: "user1/123_image.png"}

	S3Client.ObjectExistsFunc = func(ctx context.Context, req dto.S3ObjectExistsRequest) (dto.S3ObjectExistsResponse, error) {
		return dto.S3ObjectExistsResponse{Exists: true}, nil
	}

	ImageRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Image) error {
		return assert.AnError
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1, Username: "user1"})

	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	res, err := u.ConfirmUpload(ctx, *req)

	require.Equal(t, dto.ImageResponse{}, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}

func TestImageUsecaseImpl_ConfirmUpload_Fail_Send(t *testing.T) {
	cfg := config.NewConfig()
	gormDB, mockDB := newFakeDB(t)
	ImageRepository := &mock.ImageRepositoryMock{}
	ImageProducer := &mock.ImageProducerMock{}
	S3Client := &mock.S3ClientMock{}

	u := &imageusecase.ImageUsecaseImpl{
		Cfg:             cfg,
		DB:              gormDB,
		ImageRepository: ImageRepository,
		ImageProducer:   ImageProducer,
		S3Client:        S3Client,
	}

	req := &dto.ConfirmUploadImageRequest{Key: "user1/123_image.png"}

	S3Client.ObjectExistsFunc = func(ctx context.Context, req dto.S3ObjectExistsRequest) (dto.S3ObjectExistsResponse, error) {
		return dto.S3ObjectExistsResponse{Exists: true}, nil
	}

	ImageRepository.CreateFunc = func(ctx context.Context, db *gorm.DB, entityMoqParam *entity.Image) error {
		return nil
	}

	ImageProducer.SendImageUploadedFunc = func(ctx context.Context, db *gorm.DB, event *dto.ImageUploadedEvent) error {
		return assert.AnError
	}

	ctx := context.Background()
	ctx = ctxuserauth.Set(ctx, &dto.UserAuth{ID: 1, Username: "user1"})

	mockDB.ExpectBegin()
	mockDB.ExpectRollback()

	res, err := u.ConfirmUpload(ctx, *req)

	require.Equal(t, dto.ImageResponse{}, res)
	require.NotNil(t, err)
	require.ErrorIs(t, err, assert.AnError)
}
