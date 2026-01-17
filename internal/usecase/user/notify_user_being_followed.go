package user

import (
	"context"
	"fmt"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *UserUsecaseImpl) NotifyUserBeingFollowed(ctx context.Context, req *model.NotifyUserBeingFollowedRequest) error {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	followerUser := new(entity.User)

	if err := u.UserRepository.FindByID(ctx, u.DB, followerUser, req.FollowerID); err != nil {
		return errkit.AddFuncName(err)
	}

	event := &model.NotifEvent{
		UserID:  req.FollowingID,
		Message: fmt.Sprintf("%s just follow you", followerUser.Name),
	}

	if err := u.NotifProducer.SendNotif(ctx, event); err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
