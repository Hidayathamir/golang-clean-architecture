package rest

import (
	"context"

	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/S3Client.go -pkg=mock . S3Client

type S3Client interface {
	Download(ctx context.Context, bucket, key string) (string, error)
	DeleteObject(ctx context.Context, bucket, key string) (bool, error)
}

var _ S3Client = &S3ClientImpl{}

type S3ClientImpl struct {
	Config *viper.Viper
}

func NewS3Client(cfg *viper.Viper) *S3ClientImpl {
	return &S3ClientImpl{
		Config: cfg,
	}
}

func (c *S3ClientImpl) Download(ctx context.Context, bucket, key string) (string, error) {
	// TODO implement hit external rest api
	return "dummy data", nil
}

func (c *S3ClientImpl) DeleteObject(ctx context.Context, bucket, key string) (bool, error) {
	// TODO implement hit external rest api
	return true, nil
}
