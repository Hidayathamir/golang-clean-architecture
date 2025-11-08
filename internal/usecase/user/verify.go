package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *UserUsecaseImpl) Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error) {
	err := x.Validate.Struct(req)
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

	auth := new(model.Auth)
	converter.EntityUserToModelAuth(user, auth)

	return auth, nil
}
