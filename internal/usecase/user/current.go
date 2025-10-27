package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *UserUsecaseImpl) Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Current", err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByID(ctx, u.DB.WithContext(ctx), user, req.ID); err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Current", err)
	}

	res := new(model.UserResponse)
	converter.UserToResponse(user, res)

	return res, nil
}
