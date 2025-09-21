package restmwlogger

import (
	"context"
	"golang-clean-architecture/internal/gateway/rest"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ rest.S3Client = &S3ClientImpl{}

type S3ClientImpl struct {
	logger *logrus.Logger

	next rest.S3Client
}

func NewS3Client(logger *logrus.Logger, next rest.S3Client) *S3ClientImpl {
	return &S3ClientImpl{
		logger: logger,
		next:   next,
	}
}

func (c *S3ClientImpl) Download(ctx context.Context, bucket, key string) (string, error) {
	ok, err := c.next.Download(ctx, bucket, key)

	fields := logrus.Fields{
		"bucket": bucket,
		"key":    key,
		"ok":     ok,
	}
	helper.Log(ctx, fields, err)

	return ok, err
}

func (c *S3ClientImpl) DeleteObject(ctx context.Context, bucket, key string) (bool, error) {
	ok, err := c.next.DeleteObject(ctx, bucket, key)

	fields := logrus.Fields{
		"bucket": bucket,
		"key":    key,
		"ok":     ok,
	}
	helper.Log(ctx, fields, err)

	return ok, err
}
