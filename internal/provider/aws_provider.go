package provider

import (
	"context"
	"errors"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/bucketname"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

func NewAWSS3Client(cfg *config.Config) *s3.Client {
	region := cfg.GetAWSRegion()

	awsConfig, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "testtesttest", "")),
		awsconfig.WithRetryMaxAttempts(3),
	)
	if err != nil {
		logkit.Logger.Panic(err)
	}

	baseEndpoint := cfg.GetAWSBaseEndpoint()

	awsS3Client := s3.NewFromConfig(awsConfig, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(baseEndpoint)
		o.UsePathStyle = true
	})

	ensureBucketExists(awsS3Client, bucketname.Image)

	return awsS3Client
}

// ensureBucketExists creates the bucket when missing so a fresh environment
// (local MinIO) works without a manual setup step. Both binaries call
// NewAWSS3Client, so this is the single bootstrap point. S3 may not be ready
// the instant the app boots, hence the retry loop, mirroring NewDatabase.
func ensureBucketExists(client *s3.Client, bucket string) {
	const maxAttempts = 5

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		err := createBucketIfMissing(client, bucket)
		if err == nil {
			return
		}
		logkit.Logger.Warnf("s3 bucket %q attempt %d/%d failed: %v", bucket, attempt, maxAttempts, err)
		time.Sleep(1 * time.Second)
	}

	logkit.Logger.Panicf("failed to ensure s3 bucket %q", bucket)
}

func createBucketIfMissing(client *s3.Client, bucket string) error {
	ctx := context.Background()

	if _, err := client.HeadBucket(ctx, &s3.HeadBucketInput{Bucket: &bucket}); err == nil {
		return nil
	}

	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{Bucket: &bucket})
	if err == nil {
		logkit.Logger.Infof("created s3 bucket %q", bucket)
		return nil
	}

	// A concurrent boot may have created it first; that is success, not failure.
	var alreadyExists *types.BucketAlreadyExists
	var alreadyOwned *types.BucketAlreadyOwnedByYou
	if errors.As(err, &alreadyExists) || errors.As(err, &alreadyOwned) {
		return nil
	}

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && apiErr.ErrorCode() == "BucketAlreadyOwnedByYou" {
		return nil
	}

	return err
}
