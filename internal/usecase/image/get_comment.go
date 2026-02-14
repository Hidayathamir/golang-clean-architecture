package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) GetComment(ctx context.Context, req dto.GetCommentRequest) (dto.CommentResponseList, error) {
	err := x.Validate.Struct(&req)
	if err != nil {
		err = errkit.BadRequest(err)
		return dto.CommentResponseList{}, errkit.AddFuncName(err)
	}

	commentList := entity.CommentList{}
	err = u.CommentRepository.FindByImageID(ctx, u.DB, &commentList, req.ImageID)
	if err != nil {
		return dto.CommentResponseList{}, errkit.AddFuncName(err)
	}

	res := dto.CommentResponseList{}
	converter.EntityCommentListToDtoCommentResponseList(commentList, &res)

	return res, nil
}
