package userusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func (u *UserUsecaseImpl) GetProfile(ctx context.Context, req dto.GetProfileRequest) (dto.ProfileResponse, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.ProfileResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).GetProfile")
	}

	user := entity.User{}
	err = u.UserRepository.FindByID(ctx, u.DB, &user, req.ID)
	if err != nil {
		return dto.ProfileResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).GetProfile")
	}

	userStat := entity.UserStat{}
	err = u.UserStatRepository.FindByUserID(ctx, u.DB, &userStat, req.ID)
	if err != nil {
		return dto.ProfileResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).GetProfile")
	}

	res := dto.ProfileResponse{}
	converter.EntityUserAndUserStatToDtoProfileResponse(user, userStat, &res)

	return res, nil
}
