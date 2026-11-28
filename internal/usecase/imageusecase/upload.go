package imageusecase

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/s3key"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"gorm.io/gorm"
)

const presignUploadExpiry = 15 * time.Minute

func (u *ImageUsecaseImpl) PresignUpload(ctx context.Context, req dto.PresignUploadImageRequest) (dto.PresignUploadImageResponse, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.PresignUploadImageResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).PresignUpload")
	}

	s3PresignReq := dto.S3PresignUploadRequest{ExpiresIn: presignUploadExpiry}
	converter.DtoPresignUploadImageRequestToDtoS3PresignUploadRequest(ctx, req, &s3PresignReq)

	presigned, err := u.S3Client.PresignUploadImage(ctx, s3PresignReq)
	if err != nil {
		return dto.PresignUploadImageResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).PresignUpload")
	}

	res := dto.PresignUploadImageResponse{
		UploadURL: presigned.UploadURL,
		Key:       presigned.Key,
		URL:       s3key.ImageURL(u.Cfg.GetAWSBaseEndpoint(), presigned.Key),
	}

	return res, nil
}

func (u *ImageUsecaseImpl) ConfirmUpload(ctx context.Context, req dto.ConfirmUploadImageRequest) (dto.ImageResponse, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.ImageResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).ConfirmUpload")
	}

	exists, err := u.S3Client.ObjectExists(ctx, dto.S3ObjectExistsRequest{Key: req.Key})
	if err != nil {
		return dto.ImageResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).ConfirmUpload")
	}
	if !exists.Exists {
		err = errkit.SetCode(errors.New("upload not found in storage"), http.StatusBadRequest)
		return dto.ImageResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).ConfirmUpload")
	}

	image := entity.Image{
		UserID:  ctxuserauth.Get(ctx).ID,
		Caption: req.Caption,
		URL:     s3key.ImageURL(u.Cfg.GetAWSBaseEndpoint(), req.Key),
	}

	err = u.DB.Transaction(func(tx *gorm.DB) error {
		err := u.ImageRepository.Create(ctx, tx, &image)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).ConfirmUpload")
		}

		event := dto.ImageUploadedEvent{}
		converter.EntityImageToDtoImageUploadedEvent(image, &event)

		err = u.ImageProducer.SendImageUploaded(ctx, tx, &event)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).ConfirmUpload")
		}

		return nil
	})
	if err != nil {
		return dto.ImageResponse{}, err
	}

	res := dto.ImageResponse{}
	converter.EntityImageToDtoImageResponse(image, &res)

	return res, nil
}
