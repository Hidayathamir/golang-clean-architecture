package imageusecase

import (
	"context"
	"net/http"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *ImageUsecaseImpl) GetUserImages(ctx context.Context, req dto.UserImagesRequest) (dto.UserImagesResponse, error) {
	cursor, hasCursor, err := req.Page.DecodeCursor()
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.UserImagesResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetUserImages")
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
	err = u.ImageRepository.FindByUserID(ctx, u.DB, &imageList, req.UserID, cursorCreatedAt, cursorID, limit+1)
	if err != nil {
		return dto.UserImagesResponse{}, errkit.AddFuncName(err, "imageusecase.(*ImageUsecaseImpl).GetUserImages")
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
