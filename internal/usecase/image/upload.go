package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

// Upload implements ImageUsecase.
func (u *ImageUsecaseImpl) Upload(ctx context.Context, req *model.UploadImageRequest) error {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	image := &entity.Image{
		UserID:       0,
		URL:          "",
		LikeCount:    0,
		CommentCount: 0,
	}

	converter.ModelUploadImageRequestToEntityImage(req, image)

	if err := u.ImageRepository.Create(ctx, u.DB, image); err != nil {
		return errkit.AddFuncName(err)
	}

	event := new(model.ImageUploadedEvent)
	converter.EntityImageToModelImageUploadedEvent(image, event)

	if err := u.ImageUploadedProducer.Send(ctx, event); err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
