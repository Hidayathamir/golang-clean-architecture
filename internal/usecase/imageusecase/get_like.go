package imageusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) GetLike(ctx context.Context, req dto.GetLikeRequest) (dto.LikeResponseList, error) {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.LikeResponseList{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetLike")
	}

	likeList := entity.LikeList{}
	err = u.LikeRepository.FindByImageID(ctx, u.DB, &likeList, req.ImageID)
	if err != nil {
		return dto.LikeResponseList{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetLike")
	}

	res := dto.LikeResponseList{}
	converter.EntityLikeListToDtoLikeResponseList(likeList, &res)

	return res, nil
}
