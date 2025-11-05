package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *UserUsecaseImpl) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error) {
	err := u.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Verify", err)
	}

	userID, err := u.parseAccessToken(ctx, req.Token)
	if err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Verify", err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByID(ctx, u.DB.WithContext(ctx), user, userID); err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Verify", err)
	}

	return &model.Auth{ID: user.ID}, nil
}
