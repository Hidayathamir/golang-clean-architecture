package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/spf13/viper"
)

//go:generate moq -out=../../mock/S3Client.go -pkg=mock . S3Client

type S3Client interface {
	Download(ctx context.Context, req model.S3DownloadRequest) (model.S3DownloadResponse, error)
	DeleteObject(ctx context.Context, req model.S3DeleteObjectRequest) (model.S3DeleteObjectResponse, error)
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
