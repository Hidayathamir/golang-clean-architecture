package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) Upload(ctx context.Context, req *model.UploadImageRequest) (*model.ImageResponse, error) {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	s3UploadImgReq := model.S3UploadImageRequest{}
	if err := converter.ModelUploadImageRequestToModelS3UploadImageRequest(ctx, req, &s3UploadImgReq); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	url, err := u.S3Client.UploadImage(ctx, s3UploadImgReq)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	image := &entity.Image{UserID: ctxuserauth.Get(ctx).ID, URL: url}

	if err := u.ImageRepository.Create(ctx, u.DB, image); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := new(model.ImageUploadedEvent)
	converter.EntityImageToModelImageUploadedEvent(ctx, image, event)

	if err := u.ImageProducer.SendImageUploaded(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(model.ImageResponse)
	converter.EntityImageToModelImageResponse(ctx, image, res)

	return res, nil
}
