package notifusecase

import (
	"context"
	"errors"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/ctx/ctxuserauth"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func (u *NotifUsecaseImpl) MarkRead(ctx context.Context, req dto.MarkNotificationReadRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).MarkRead")
	}
	userAuth := ctxuserauth.Get(ctx)
	affected, err := u.NotificationRepository.MarkReadByID(ctx, u.DB, userAuth.ID, req.NotificationID)
	if err != nil {
		return errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).MarkRead")
	}
	if affected == 0 {
		err := errors.New("notification not found")
		err = errkit.SetCode(err, http.StatusNotFound)
		return errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).MarkRead")
	}
	return nil
}

func (u *NotifUsecaseImpl) MarkAllRead(ctx context.Context) error {
	userAuth := ctxuserauth.Get(ctx)
	_, err := u.NotificationRepository.MarkAllReadByUserID(ctx, u.DB, userAuth.ID)
	if err != nil {
		return errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).MarkAllRead")
	}
	return nil
}

func (u *NotifUsecaseImpl) CountUnread(ctx context.Context) (dto.UnreadCountResponse, error) {
	userAuth := ctxuserauth.Get(ctx)
	count, err := u.NotificationRepository.CountUnreadByUserID(ctx, u.DB, userAuth.ID)
	if err != nil {
		return dto.UnreadCountResponse{}, errkit.AddFuncName(err, "notifusecase.(*NotifUsecaseImpl).CountUnread")
	}
	return dto.UnreadCountResponse{UnreadCount: count}, nil
}
