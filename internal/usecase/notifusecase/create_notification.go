package notifusecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
)

// createAndPush persists the notification then best-effort pushes it to any
// live connection. Self-actions are dropped here so every call site gets the
// rule for free.
func (u *NotifUsecaseImpl) createAndPush(ctx context.Context, req dto.NotifyRequest) error {
	if req.ActorID == req.UserID {
		return nil
	}
	notification := entity.Notification{}
	notification.UserID = req.UserID
	notification.ActorID = req.ActorID
	notification.Type = req.Type
	notification.ImageID = req.ImageID
	if !notification.IsValid() {
		err := errors.New("unknown notification type")
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).createAndPush")
	}
	err := u.NotificationRepository.Create(ctx, u.DB, &notification)
	if err != nil {
		return errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).createAndPush")
	}
	actor := entity.User{}
	err = u.UserRepository.FindByID(ctx, u.DB, &actor, req.ActorID)
	if err != nil {
		return errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).createAndPush")
	}
	res := dto.NotificationResponse{}
	res.ID = notification.ID
	res.ActorID = notification.ActorID
	res.ActorName = actor.Name
	res.Type = notification.Type
	res.ImageID = notification.ImageID
	res.IsRead = notification.IsRead
	res.CreatedAt = notification.CreatedAt
	err = u.NotificationPublisher.Publish(ctx, notification.UserID, res)
	if err != nil {
		return errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).createAndPush")
	}
	logkit.Logger.WithContext(ctx).WithField("notification", res).Debug("notification published")
	return nil
}
