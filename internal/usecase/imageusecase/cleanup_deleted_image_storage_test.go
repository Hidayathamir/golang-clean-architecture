package imageusecase_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/imageusecase"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/bucketname"
	"github.com/stretchr/testify/require"
)

func TestImageUsecaseImpl_CleanupDeletedImageStorage_Success(t *testing.T) {
	cfg := config.NewConfig()
	s3Client := &mock.S3ClientMock{}

	u := &imageusecase.ImageUsecaseImpl{
		Cfg:      cfg,
		S3Client: s3Client,
	}

	url := cfg.GetAWSBaseEndpoint() + "/" + bucketname.Image + "/user/123_file.jpg"

	var got dto.S3DeleteObjectRequest
	s3Client.DeleteObjectFunc = func(ctx context.Context, req dto.S3DeleteObjectRequest) (dto.S3DeleteObjectResponse, error) {
		got = req
		return dto.S3DeleteObjectResponse{Deleted: true}, nil
	}

	err := u.CleanupDeletedImageStorage(context.Background(), dto.CleanupDeletedImageStorageRequest{ID: 5, URL: url})

	require.Nil(t, err)
	require.Equal(t, bucketname.Image, got.Bucket)
	require.Equal(t, "user/123_file.jpg", got.Key)
}

func TestImageUsecaseImpl_CleanupDeletedImageStorage_Fail_DeleteObject(t *testing.T) {
	cfg := config.NewConfig()
	s3Client := &mock.S3ClientMock{}

	u := &imageusecase.ImageUsecaseImpl{
		Cfg:      cfg,
		S3Client: s3Client,
	}

	s3Client.DeleteObjectFunc = func(ctx context.Context, req dto.S3DeleteObjectRequest) (dto.S3DeleteObjectResponse, error) {
		return dto.S3DeleteObjectResponse{}, errDeleteS3
	}

	err := u.CleanupDeletedImageStorage(context.Background(), dto.CleanupDeletedImageStorageRequest{ID: 5, URL: "http://x/image/k.jpg"})

	require.ErrorIs(t, err, errDeleteS3)
}
