package imageusecase

import (
	"context"
	"net/http"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *ImageUsecaseImpl) GetFeed(ctx context.Context, req dto.FeedRequest) (dto.FeedResponse, error) {
	userAuth := ctxuserauth.Get(ctx)

	cursor, hasCursor, err := req.Page.DecodeCursor()
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.FeedResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetFeed")
	}

	var cursorCreatedAt time.Time
	var cursorID int64
	if hasCursor {
		cursorCreatedAt = cursor.CreatedAt
		cursorID = cursor.ID
	}

	limit := req.Page.NormalizedLimit()

	// fetch one extra row to know whether a further page exists.
	imageList := entity.ImageList{}
	err = u.ImageRepository.FindFeedByFollowerID(ctx, u.DB, &imageList, userAuth.ID, cursorCreatedAt, cursorID, limit+1)
	if err != nil {
		return dto.FeedResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetFeed")
	}

	hasMore := len(imageList) > limit
	if hasMore {
		imageList = imageList[:limit]
	}

	res := dto.ImageResponseList{}
	converter.EntityImageListToDtoImageResponseList(imageList, &res)

	var lastCreatedAt time.Time
	var lastID int64
	if len(imageList) > 0 {
		last := imageList[len(imageList)-1]
		lastCreatedAt = last.CreatedAt
		lastID = last.ID
	}

	return dto.PageFrom(res, hasMore, lastCreatedAt, lastID), nil
}
