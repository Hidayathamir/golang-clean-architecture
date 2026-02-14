package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) GetLike(ctx context.Context, req dto.GetLikeRequest) (dto.LikeResponseList, error) {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		return dto.LikeResponseList{}, errkit.AddFuncName(err)
	}

	likeList := entity.LikeList{}
	err = u.LikeRepository.FindByImageID(ctx, u.DB, &likeList, req.ImageID)
	if err != nil {
		return dto.LikeResponseList{}, errkit.AddFuncName(err)
	}

	res := dto.LikeResponseList{}
	converter.EntityLikeListToDtoLikeResponseList(likeList, &res)

	return res, nil
}
