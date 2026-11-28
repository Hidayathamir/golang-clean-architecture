package storage

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/bucketname"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/s3key"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:generate moq -out=../../mock/MockClientS3.go -pkg=mock . S3Client

type S3Client interface {
	UploadImage(ctx context.Context, req dto.S3UploadImageRequest) (url string, err error)
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
