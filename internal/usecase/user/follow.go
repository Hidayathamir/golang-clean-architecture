package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *UserUsecaseImpl) Follow(ctx context.Context, req *model.FollowUserRequest) error {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	follow := new(entity.Follow)
	converter.ModelFollowUserRequestToEntityFollow(ctx, req, follow)

	if err := u.FollowRepository.Create(ctx, u.DB, follow); err != nil {
		return errkit.AddFuncName(err)
	}

	event := new(model.UserFollowedEvent)
	converter.EntityFollowToModelUserFollowedEvent(ctx, follow, event)

	if err := u.UserProducer.SendUserFollowed(ctx, event); err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
