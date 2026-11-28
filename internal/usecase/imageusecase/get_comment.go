package imageusecase

import (
	"context"
	"net/http"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func (u *ImageUsecaseImpl) GetComment(ctx context.Context, req dto.GetCommentRequest) (dto.CommentsResponse, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.CommentsResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetComment")
	}

	cursor, hasCursor, err := req.Page.DecodeCursor()
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.CommentsResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetComment")
	}

	var cursorCreatedAt time.Time
	var cursorID int64
	if hasCursor {
		cursorCreatedAt = cursor.CreatedAt
		cursorID = cursor.ID
	}

	limit := req.Page.NormalizedLimit()

	// fetch one extra row to know whether a further page exists.
	commentList := entity.CommentList{}
	err = u.CommentRepository.FindByImageID(ctx, u.DB, &commentList, req.ImageID, cursorCreatedAt, cursorID, limit+1)
	if err != nil {
		return dto.CommentsResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetComment")
	}

	hasMore := len(commentList) > limit
	if hasMore {
		commentList = commentList[:limit]
	}

	res := dto.CommentResponseList{}
	converter.EntityCommentListToDtoCommentResponseList(commentList, &res)

	var lastCreatedAt time.Time
	var lastID int64
	if len(commentList) > 0 {
		last := commentList[len(commentList)-1]
		lastCreatedAt = last.CreatedAt
		lastID = last.ID
	}

	return dto.PageFrom(res, hasMore, lastCreatedAt, lastID), nil
}
