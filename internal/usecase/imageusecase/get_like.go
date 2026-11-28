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

func (u *ImageUsecaseImpl) GetLike(ctx context.Context, req dto.GetLikeRequest) (dto.LikesResponse, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.LikesResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetLike")
	}

	cursor, hasCursor, err := req.Page.DecodeCursor()
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.LikesResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetLike")
	}

	var cursorCreatedAt time.Time
	var cursorID int64
	if hasCursor {
		cursorCreatedAt = cursor.CreatedAt
		cursorID = cursor.ID
	}

	limit := req.Page.NormalizedLimit()

	// fetch one extra row to know whether a further page exists.
	likeList := entity.LikeList{}
	err = u.LikeRepository.FindByImageID(ctx, u.DB, &likeList, req.ImageID, cursorCreatedAt, cursorID, limit+1)
	if err != nil {
		return dto.LikesResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetLike")
	}

	hasMore := len(likeList) > limit
	if hasMore {
		likeList = likeList[:limit]
	}

	res := dto.LikeResponseList{}
	converter.EntityLikeListToDtoLikeResponseList(likeList, &res)

	var lastCreatedAt time.Time
	var lastID int64
	if len(likeList) > 0 {
		last := likeList[len(likeList)-1]
		lastCreatedAt = last.CreatedAt
		lastID = last.ID
	}

	return dto.PageFrom(res, hasMore, lastCreatedAt, lastID), nil
}
