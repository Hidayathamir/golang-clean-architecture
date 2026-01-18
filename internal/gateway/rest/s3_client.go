package rest

import (
	"context"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/bucketname"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/MockClientS3.go -pkg=mock . S3Client

type S3Client interface {
	UploadImage(ctx context.Context, req model.S3UploadImageRequest) (url string, err error)
	Download(ctx context.Context, req model.S3DownloadRequest) (model.S3DownloadResponse, error)
	DeleteObject(ctx context.Context, req model.S3DeleteObjectRequest) (model.S3DeleteObjectResponse, error)
}

var _ S3Client = &S3ClientImpl{}

type S3ClientImpl struct {
	Config      *viper.Viper
	AWSS3Client *s3.Client
}

func NewS3Client(cfg *viper.Viper, awsS3Client *s3.Client) *S3ClientImpl {
	return &S3ClientImpl{
		Config:      cfg,
		AWSS3Client: awsS3Client,
	}
}

func (c *S3ClientImpl) UploadImage(ctx context.Context, req model.S3UploadImageRequest) (url string, err error) {
	_, err = c.AWSS3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucketname.Image,
		Key:    &req.Key,
		Body:   req.Body,
	})
	if err != nil {
		return "", errkit.AddFuncName(err)
	}

	baseEndpoint := c.Config.GetString(configkey.AWSBaseEndpoint)
	url = fmt.Sprintf("%s/%s/%s", baseEndpoint, bucketname.Image, req.Key)

	return url, nil
}

func (c *S3ClientImpl) Download(ctx context.Context, req model.S3DownloadRequest) (model.S3DownloadResponse, error) {
	// TODO implement hit external rest api
	return model.S3DownloadResponse{
		Data: "dummy data",
	}, nil
}

func (c *S3ClientImpl) DeleteObject(ctx context.Context, req model.S3DeleteObjectRequest) (model.S3DeleteObjectResponse, error) {
	// TODO implement hit external rest api
	return model.S3DeleteObjectResponse{
		Deleted: true,
	}, nil
}
