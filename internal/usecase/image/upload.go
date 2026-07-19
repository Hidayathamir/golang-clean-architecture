package image

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) Upload(ctx context.Context, req dto.UploadImageRequest) (dto.ImageResponse, error) {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.ImageResponse{}, errkit.AddFuncName(err, "image.(*ImageUsecaseImpl).Upload")
	}

	s3UploadImgReq := dto.S3UploadImageRequest{}
	err = converter.DtoUploadImageRequestToDtoS3UploadImageRequest(ctx, req, &s3UploadImgReq)
	if err != nil {
		return dto.ImageResponse{}, errkit.AddFuncName(err, "image.(*ImageUsecaseImpl).Upload")
	}

	url, err := u.S3Client.UploadImage(ctx, s3UploadImgReq)
	if err != nil {
		return dto.ImageResponse{}, errkit.AddFuncName(err, "image.(*ImageUsecaseImpl).Upload")
	}

	image := entity.Image{UserID: ctxuserauth.Get(ctx).ID, Caption: req.Caption, URL: url}

	err = u.ImageRepository.Create(ctx, u.DB, &image)
	if err != nil {
		return dto.ImageResponse{}, errkit.AddFuncName(err, "image.(*ImageUsecaseImpl).Upload")
	}

	event := dto.ImageUploadedEvent{}
	converter.EntityImageToDtoImageUploadedEvent(image, &event)

	err = u.ImageProducer.SendImageUploaded(ctx, &event)
	if err != nil {
		return dto.ImageResponse{}, errkit.AddFuncName(err, "image.(*ImageUsecaseImpl).Upload")
	}

	res := dto.ImageResponse{}
	converter.EntityImageToDtoImageResponse(image, &res)

	return res, nil
}
