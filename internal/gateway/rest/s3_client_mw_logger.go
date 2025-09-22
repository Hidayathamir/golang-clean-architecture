package rest

import (
	"context"
	"golang-clean-architecture/pkg/helper"

	"github.com/sirupsen/logrus"
)

var _ S3Client = &S3ClientMwLogger{}

type S3ClientMwLogger struct {
	logger *logrus.Logger

	next S3Client
}

func NewS3ClientMwLogger(logger *logrus.Logger, next S3Client) *S3ClientMwLogger {
	return &S3ClientMwLogger{
		logger: logger,
		next:   next,
	}
}

func (c *S3ClientMwLogger) Download(ctx context.Context, bucket, key string) (string, error) {
	ok, err := c.next.Download(ctx, bucket, key)

	fields := logrus.Fields{
		"bucket": bucket,
		"key":    key,
		"ok":     ok,
	}
	helper.Log(ctx, fields, err)

	return ok, err
}

func (c *S3ClientMwLogger) DeleteObject(ctx context.Context, bucket, key string) (bool, error) {
	ok, err := c.next.DeleteObject(ctx, bucket, key)

	fields := logrus.Fields{
		"bucket": bucket,
		"key":    key,
		"ok":     ok,
	}
	helper.Log(ctx, fields, err)

	return ok, err
}
