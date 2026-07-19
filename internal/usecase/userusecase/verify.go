package userusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func (u *UserUsecaseImpl) Verify(ctx context.Context, req dto.VerifyUserRequest) (dto.UserAuth, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.UserAuth{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Verify")
	}

	userID, err := u.parseAccessToken(ctx, req.Token)
	if err != nil {
		return dto.UserAuth{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Verify")
	}

	cachedUser, err := u.UserCache.Get(ctx, userID)
	if err == nil && cachedUser != nil {
		userAuth := dto.UserAuth{}
		converter.EntityUserToDtoUserAuth(*cachedUser, &userAuth)
		return userAuth, nil
	}

	user := entity.User{}
	err = u.UserRepository.FindByID(ctx, u.DB, &user, userID)
	if err != nil {
		return dto.UserAuth{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Verify")
	}

	err = u.UserCache.Set(ctx, &user)
	logkit.LogIfErr(err)

	userAuth := dto.UserAuth{}
	converter.EntityUserToDtoUserAuth(user, &userAuth)

	return userAuth, nil
}
