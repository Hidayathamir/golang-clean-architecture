package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) Comment(ctx context.Context, req dto.CommentImageRequest) error {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	comment := entity.Comment{}
	converter.DtoCommentImageRequestToEntityComment(ctx, req, &comment)

	err = u.CommentRepository.Create(ctx, u.DB, &comment)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	event := dto.ImageCommentedEvent{}
	converter.EntityCommentToDtoImageCommentedEvent(comment, &event)

	err = u.ImageProducer.SendImageCommented(ctx, &event)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
