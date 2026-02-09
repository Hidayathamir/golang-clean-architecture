package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) GetComment(ctx context.Context, req *dto.GetCommentRequest) (*dto.CommentResponseList, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	commentList := new(entity.CommentList)
	if err := u.CommentRepository.FindByImageID(ctx, u.DB, commentList, req.ImageID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(dto.CommentResponseList)
	converter.EntityCommentListToDtoCommentResponseList(ctx, commentList, res)

	return res, nil
}
