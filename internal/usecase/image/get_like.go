package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) GetLike(ctx context.Context, req *model.GetLikeRequest) (*model.LikeResponseList, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	likeList := new(entity.LikeList)
	if err := u.LikeRepository.FindByImageID(ctx, u.DB, likeList, req.ImageID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(model.LikeResponseList)
	converter.EntityLikeListToModelLikeResponseList(ctx, likeList, res)

	return res, nil
}
