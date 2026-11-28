package userusecase

import (
	"context"
	"net/http"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *UserUsecaseImpl) GetFollowers(ctx context.Context, req dto.FollowersRequest) (dto.FollowersResponse, error) {
	cursor, hasCursor, err := req.Page.DecodeCursor()
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.FollowersResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).GetFollowers")
	}

	var cursorCreatedAt time.Time
	var cursorID int64
	if hasCursor {
		cursorCreatedAt = cursor.CreatedAt
		cursorID = cursor.ID
	}

	limit := req.Page.NormalizedLimit()

	// fetch one extra row to know whether a further page exists.
	followList := entity.FollowList{}
	err = u.FollowRepository.FindFollowersByFollowingID(ctx, u.DB, &followList, req.UserID, cursorCreatedAt, cursorID, limit+1)
	if err != nil {
		return dto.FollowersResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).GetFollowers")
	}

	hasMore := len(followList) > limit
	if hasMore {
		followList = followList[:limit]
	}

	var lastCreatedAt time.Time
	var lastID int64
	if len(followList) > 0 {
		last := followList[len(followList)-1]
		lastCreatedAt = last.CreatedAt
		lastID = last.ID
	}

	followerIDs := make([]int64, 0, len(followList))
	for _, follow := range followList {
		followerIDs = append(followerIDs, follow.FollowerID)
	}

	userList := entity.UserList{}
	if len(followerIDs) > 0 {
		err = u.UserRepository.FindByIDs(ctx, u.DB, &userList, followerIDs)
		if err != nil {
			return dto.FollowersResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).GetFollowers")
		}
	}

	// FindByIDs does not preserve input order; reorder to match the follow page.
	userByID := make(map[int64]entity.User, len(userList))
	for _, user := range userList {
		userByID[user.ID] = user
	}
	ordered := entity.UserList{}
	for _, id := range followerIDs {
		if user, ok := userByID[id]; ok {
			ordered = append(ordered, user)
		}
	}

	res := dto.UserResponseList{}
	converter.EntityUserListToDtoUserResponseList(ordered, &res)

	return dto.PageFrom(res, hasMore, lastCreatedAt, lastID), nil
}
