package user

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/pkg/errkit"
)

func (u *UserUsecaseImpl) Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindById(ctx, tx, user, req.ID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.UserToResponse(user), nil
}
