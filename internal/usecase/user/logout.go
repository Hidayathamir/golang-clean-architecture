package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *UserUsecaseImpl) Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return false, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(ctx, tx, user, req.ID); err != nil {
		return false, errkit.AddFuncName(err)
	}

	user.Token = ""

	if err := u.UserRepository.Update(ctx, tx, user); err != nil {
		return false, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return false, errkit.AddFuncName(err)
	}

	if _, err := u.S3Client.DeleteObject(ctx, "user-bucket", user.ID); err != nil {
		return false, errkit.AddFuncName(err)
	}

	event := converter.UserToEvent(user)
	if err := u.UserProducer.Send(ctx, event); err != nil {
		return false, errkit.AddFuncName(err)
	}

	return true, nil
}
