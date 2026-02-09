package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) Comment(ctx context.Context, req *dto.CommentImageRequest) error {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	comment := new(entity.Comment)
	converter.DtoCommentImageRequestToEntityComment(ctx, req, comment)

	if err := u.CommentRepository.Create(ctx, u.DB, comment); err != nil {
		return errkit.AddFuncName(err)
	}

	event := new(dto.ImageCommentedEvent)
	converter.EntityCommentToDtoImageCommentedEvent(ctx, comment, event)

	if err := u.ImageProducer.SendImageCommented(ctx, event); err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
