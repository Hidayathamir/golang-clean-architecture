package notifusecase

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

func (u *NotifUsecaseImpl) List(ctx context.Context, req dto.NotificationsRequest) (dto.NotificationsResponse, error) {
	userAuth := ctxuserauth.Get(ctx)
	cursor, hasCursor, err := req.Page.DecodeCursor()
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.NotificationsResponse{}, errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).List")
	}
	var cursorCreatedAt time.Time
	var cursorID int64
	if hasCursor {
		cursorCreatedAt = cursor.CreatedAt
		cursorID = cursor.ID
	}
	limit := req.Page.NormalizedLimit()
	notificationList := entity.NotificationList{}
	err = u.NotificationRepository.FindByUserID(ctx, u.DB, &notificationList, userAuth.ID, cursorCreatedAt, cursorID, limit+1)
	if err != nil {
		return dto.NotificationsResponse{}, errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).List")
	}
	hasMore := len(notificationList) > limit
	if hasMore {
		notificationList = notificationList[:limit]
	}
	res := dto.NotificationResponseList{}
	converter.EntityNotificationListToDtoNotificationResponseList(notificationList, &res)
	err = u.attachActorNames(ctx, notificationList, &res)
	if err != nil {
		return dto.NotificationsResponse{}, errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).List")
	}
	var lastCreatedAt time.Time
	var lastID int64
	if len(notificationList) > 0 {
		last := notificationList[len(notificationList)-1]
		lastCreatedAt = last.CreatedAt
		lastID = last.ID
	}
	return dto.PageFrom(res, hasMore, lastCreatedAt, lastID), nil
}

// attachActorNames batch-loads the distinct actors so the list is one extra
// query regardless of page size.
func (u *NotifUsecaseImpl) attachActorNames(ctx context.Context, notificationList entity.NotificationList, res *dto.NotificationResponseList) error {
	seen := map[int64]struct{}{}
	ids := make([]int64, 0, len(notificationList))
	for i := range notificationList {
		actorID := notificationList[i].ActorID
		if _, ok := seen[actorID]; ok {
			continue
		}
		seen[actorID] = struct{}{}
		ids = append(ids, actorID)
	}
	if len(ids) == 0 {
		return nil
	}
	userList := entity.UserList{}
	err := u.UserRepository.FindByIDs(ctx, u.DB, &userList, ids)
	if err != nil {
		return errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).attachActorNames")
	}
	nameByID := make(map[int64]string, len(userList))
	for i := range userList {
		nameByID[userList[i].ID] = userList[i].Name
	}
	for i := range *res {
		(*res)[i].ActorName = nameByID[(*res)[i].ActorID]
	}
	return nil
}
