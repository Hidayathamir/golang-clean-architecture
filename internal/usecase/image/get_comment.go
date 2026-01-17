package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ImageUsecaseImpl) GetComment(ctx context.Context, req *model.GetCommentRequest) (*model.CommentResponseList, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	commentList := new(entity.CommentList)
	if err := u.CommentRepository.FindByImageID(ctx, u.DB, commentList, req.ImageID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(model.CommentResponseList)
	converter.EntityCommentListToModelCommentResponseList(ctx, commentList, res)

	return res, nil
}
