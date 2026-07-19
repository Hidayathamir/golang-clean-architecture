package imageusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"gorm.io/gorm"
)

func (u *ImageUsecaseImpl) Upload(ctx context.Context, req dto.UploadImageRequest) (dto.ImageResponse, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.ImageResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Upload")
	}

	s3UploadImgReq := dto.S3UploadImageRequest{}
	err = converter.DtoUploadImageRequestToDtoS3UploadImageRequest(ctx, req, &s3UploadImgReq)
	if err != nil {
		return dto.ImageResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Upload")
	}

	url, err := u.S3Client.UploadImage(ctx, s3UploadImgReq)
	if err != nil {
		return dto.ImageResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Upload")
	}

	image := entity.Image{UserID: ctxuserauth.Get(ctx).ID, Caption: req.Caption, URL: url}

	err = u.DB.Transaction(func(tx *gorm.DB) error {
		err := u.ImageRepository.Create(ctx, tx, &image)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Upload")
		}

		event := dto.ImageUploadedEvent{}
		converter.EntityImageToDtoImageUploadedEvent(image, &event)

		err = u.ImageProducer.SendImageUploaded(ctx, tx, &event)
		if err != nil {
			return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Upload")
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
