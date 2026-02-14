package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) GetImage(ctx context.Context, req dto.GetImageRequest) (dto.ImageResponse, error) {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		return dto.ImageResponse{}, errkit.AddFuncName(err)
	}

	image := entity.Image{}
	err = u.ImageRepository.FindByID(ctx, u.DB, &image, req.ID)
	if err != nil {
		return dto.ImageResponse{}, errkit.AddFuncName(err)
	}

	res := dto.ImageResponse{}
	converter.EntityImageToDtoImageResponse(image, &res)

	return res, nil
}
