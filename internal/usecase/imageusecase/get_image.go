package imageusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func (u *ImageUsecaseImpl) GetImage(ctx context.Context, req dto.GetImageRequest) (dto.ImageResponse, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.ImageResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetImage")
	}

	image := entity.Image{}
	err = u.ImageRepository.FindByID(ctx, u.DB, &image, req.ID)
	if err != nil {
		return dto.ImageResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetImage")
	}

	res := dto.ImageResponse{}
	converter.EntityImageToDtoImageResponse(image, &res)

	return res, nil
}
