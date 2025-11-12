package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecaseImpl) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByUsername(ctx, u.DB.WithContext(ctx), user, req.Username); err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	if _, err := u.SlackClient.IsConnected(ctx, model.SlackIsConnectedRequest{}); err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	event := new(model.UserEvent)
	converter.EntityUserToModelUserEvent(user, event)
	if err := u.UserProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	token, err := u.signAccessToken(ctx, user.ID)
	if err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	res := new(model.UserResponse)
	converter.EntityUserToModelUserResponse(user, res)
	res.Token = token

	return res, nil
}
