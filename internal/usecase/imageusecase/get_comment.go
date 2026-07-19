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

func (u *ImageUsecaseImpl) GetComment(ctx context.Context, req dto.GetCommentRequest) (dto.CommentResponseList, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.CommentResponseList{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetComment")
	}

	commentList := entity.CommentList{}
	err = u.CommentRepository.FindByImageID(ctx, u.DB, &commentList, req.ImageID)
	if err != nil {
		return dto.CommentResponseList{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetComment")
	}

	res := dto.CommentResponseList{}
	converter.EntityCommentListToDtoCommentResponseList(commentList, &res)

	return res, nil
}
