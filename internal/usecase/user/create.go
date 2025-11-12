package user

import (
	"context"
	"errors"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecaseImpl) Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	err := x.Validate.Struct(req)
	if err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	total, err := u.UserRepository.CountByUsername(ctx, u.DB.WithContext(ctx), req.Username)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if total > 0 {
		err = errors.New("user already exists")
		err = errkit.Conflict(err)
		return nil, errkit.AddFuncName(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	user := new(entity.User)
	converter.ModelRegisterUserRequestToEntityUser(req, user, string(password))

	if err := u.UserRepository.Create(ctx, u.DB.WithContext(ctx), user); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := new(model.UserEvent)
	converter.EntityUserToModelUserEvent(user, event)
	if err = u.UserProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(model.UserResponse)
	converter.EntityUserToModelUserResponse(user, res)

	return res, nil
}
