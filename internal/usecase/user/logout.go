package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *UserUsecaseImpl) Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return false, errkit.AddFuncName("user.(*UserUsecaseImpl).Logout", err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByID(ctx, u.DB.WithContext(ctx), user, req.ID); err != nil {
		return false, errkit.AddFuncName("user.(*UserUsecaseImpl).Logout", err)
	}

	if _, err := u.S3Client.DeleteObject(ctx, "user-bucket", user.ID); err != nil {
		return false, errkit.AddFuncName("user.(*UserUsecaseImpl).Logout", err)
	}

	event := new(model.UserEvent)
	converter.UserToEvent(user, event)
	if err := u.UserProducer.Send(ctx, event); err != nil {
		return false, errkit.AddFuncName("user.(*UserUsecaseImpl).Logout", err)
	}

	return true, nil
}
