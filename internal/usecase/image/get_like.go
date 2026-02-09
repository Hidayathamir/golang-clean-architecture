package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) GetLike(ctx context.Context, req *dto.GetLikeRequest) (*dto.LikeResponseList, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	likeList := new(entity.LikeList)
	if err := u.LikeRepository.FindByImageID(ctx, u.DB, likeList, req.ImageID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(dto.LikeResponseList)
	converter.EntityLikeListToDtoLikeResponseList(ctx, likeList, res)

	return res, nil
}
