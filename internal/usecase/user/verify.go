package user

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/errkit"
)

func (u *UserUsecaseImpl) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := u.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByToken(ctx, tx, user, req.Token); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return &model.Auth{ID: user.ID}, nil
}
