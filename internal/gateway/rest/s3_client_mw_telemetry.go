package rest

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/pkg/telemetry"
)

var _ S3Client = &S3ClientMwTelemetry{}

type S3ClientMwTelemetry struct {
	Next S3Client
}

func NewS3ClientMwTelemetry(next S3Client) *S3ClientMwTelemetry {
	return &S3ClientMwTelemetry{
		Next: next,
	}
}

func (c *S3ClientMwTelemetry) Download(ctx context.Context, bucket, key string) (string, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	path, err := c.Next.Download(ctx, bucket, key)
	telemetry.RecordError(span, err)

	return path, err
}

func (c *S3ClientMwTelemetry) DeleteObject(ctx context.Context, bucket, key string) (bool, error) {
	ctx, span := telemetry.Start(ctx)
	defer span.End()

	ok, err := c.Next.DeleteObject(ctx, bucket, key)
	telemetry.RecordError(span, err)

	return ok, err
}
