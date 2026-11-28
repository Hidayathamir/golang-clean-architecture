package storage

import (
	"context"
	"errors"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/bucketname"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/s3key"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
)

//go:generate moq -out=../../mock/MockClientS3.go -pkg=mock . S3Client

type S3Client interface {
	UploadImage(ctx context.Context, req dto.S3UploadImageRequest) (url string, err error)
	PresignUploadImage(ctx context.Context, req dto.S3PresignUploadRequest) (dto.S3PresignUploadResponse, error)
	ObjectExists(ctx context.Context, req dto.S3ObjectExistsRequest) (dto.S3ObjectExistsResponse, error)
	DeleteObject(ctx context.Context, req dto.S3DeleteObjectRequest) (dto.S3DeleteObjectResponse, error)
}

var _ S3Client = &S3ClientImpl{}

type S3ClientImpl struct {
	Cfg         *config.Config
	AWSS3Client *s3.Client
}

func NewS3Client(cfg *config.Config, awsS3Client *s3.Client) *S3ClientImpl {
	return &S3ClientImpl{
		Cfg:         cfg,
		AWSS3Client: awsS3Client,
	}
}

func (c *S3ClientImpl) UploadImage(ctx context.Context, req dto.S3UploadImageRequest) (url string, err error) {
	_, err = c.AWSS3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketname.Image,
		Key:    &req.Key,
		Body:   req.Body,
	})
	if err != nil {
		return "", errkit.AddFuncName(err, "storage.(*S3ClientImpl).UploadImage")
	}

	baseEndpoint := c.Cfg.GetAWSBaseEndpoint()
	url = s3key.ImageURL(baseEndpoint, req.Key)

	return url, nil
}

func (c *S3ClientImpl) PresignUploadImage(ctx context.Context, req dto.S3PresignUploadRequest) (dto.S3PresignUploadResponse, error) {
	presigner := s3.NewPresignClient(c.AWSS3Client)

	presigned, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      &bucketname.Image,
		Key:         &req.Key,
		ContentType: &req.ContentType,
	}, s3.WithPresignExpires(req.ExpiresIn))
	if err != nil {
		return dto.S3PresignUploadResponse{}, errkit.AddFuncName(err, "storage.(*S3ClientImpl).PresignUploadImage")
	}

	return dto.S3PresignUploadResponse{UploadURL: presigned.URL, Key: req.Key}, nil
}

func (c *S3ClientImpl) ObjectExists(ctx context.Context, req dto.S3ObjectExistsRequest) (dto.S3ObjectExistsResponse, error) {
	_, err := c.AWSS3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucketname.Image,
		Key:    &req.Key,
	})
	if err == nil {
		return dto.S3ObjectExistsResponse{Exists: true}, nil
	}

	// A missing object is a normal, expected outcome — not an error. The caller
	// uses this to reject a confirm for an upload that never landed.
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NotFound" {
		return dto.S3ObjectExistsResponse{Exists: false}, nil
	}

	return dto.S3ObjectExistsResponse{}, errkit.AddFuncName(err, "storage.(*S3ClientImpl).ObjectExists")
}

func (c *S3ClientImpl) DeleteObject(ctx context.Context, req dto.S3DeleteObjectRequest) (dto.S3DeleteObjectResponse, error) {
	_, err := c.AWSS3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &req.Bucket,
		Key:    &req.Key,
	})
	if err != nil {
		return dto.S3DeleteObjectResponse{}, errkit.AddFuncName(err, "storage.(*S3ClientImpl).DeleteObject")
	}

	return dto.S3DeleteObjectResponse{Deleted: true}, nil
}
