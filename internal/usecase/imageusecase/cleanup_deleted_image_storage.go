package imageusecase

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/bucketname"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/s3key"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *ImageUsecaseImpl) CleanupDeletedImageStorage(ctx context.Context, req dto.CleanupDeletedImageStorageRequest) error {
	s3DeleteReq := dto.S3DeleteObjectRequest{
		Bucket: bucketname.Image,
		Key:    s3key.KeyFromImageURL(u.Cfg.GetAWSBaseEndpoint(), req.URL),
	}

	_, err := u.S3Client.DeleteObject(ctx, s3DeleteReq)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).CleanupDeletedImageStorage")
	}

	return nil
}
