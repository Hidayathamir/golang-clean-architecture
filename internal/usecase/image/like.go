package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) Like(ctx context.Context, req *dto.LikeImageRequest) error {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	like := new(entity.Like)
	converter.DtoLikeImageRequestToEntityLike(ctx, req, like)

	if err := u.LikeRepository.Create(ctx, u.DB, like); err != nil {
		return errkit.AddFuncName(err)
	}

	event := new(dto.ImageLikedEvent)
	converter.EntityLikeToDtoImageLikedEvent(ctx, like, event)

	if err := u.ImageProducer.SendImageLiked(ctx, event); err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
