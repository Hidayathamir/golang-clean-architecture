package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) GetImage(ctx context.Context, req *model.GetImageRequest) (*model.ImageResponse, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	image := new(entity.Image)
	if err := u.ImageRepository.FindByID(ctx, u.DB, image, req.ID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(model.ImageResponse)
	converter.EntityImageToModelImageResponse(ctx, image, res)

	return res, nil
}
