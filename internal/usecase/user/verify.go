package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *UserUsecaseImpl) Verify(ctx context.Context, req *dto.VerifyUserRequest) (*dto.UserAuth, error) {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	userID, err := u.parseAccessToken(ctx, req.Token)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	cachedUser, err := u.UserCache.Get(ctx, userID)
	if err == nil && cachedUser != nil {
		userAuth := new(dto.UserAuth)
		converter.EntityUserToDtoUserAuth(cachedUser, userAuth)
		return userAuth, nil
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByID(ctx, u.DB.WithContext(ctx), user, userID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	err = u.UserCache.Set(ctx, user)
	x.LogIfErr(err)

	userAuth := new(dto.UserAuth)
	converter.EntityUserToDtoUserAuth(user, userAuth)

	return userAuth, nil
}
