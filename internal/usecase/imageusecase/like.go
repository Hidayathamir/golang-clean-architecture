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

func (u *ImageUsecaseImpl) Like(ctx context.Context, req dto.LikeImageRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Like")
	}

	like := entity.Like{}
	converter.DtoLikeImageRequestToEntityLike(ctx, req, &like)

	err = u.LikeRepository.Create(ctx, u.DB, &like)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Like")
	}

	event := dto.ImageLikedEvent{}
	converter.EntityLikeToDtoImageLikedEvent(like, &event)

	err = u.ImageProducer.SendImageLiked(ctx, &event)
	if err != nil {
		return errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).Like")
	}

	return nil
}
