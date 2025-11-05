package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecaseImpl) Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByID(ctx, u.DB.WithContext(ctx), user, req.ID); err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		err = errkit.Unauthorized(err)
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	if _, err := u.SlackClient.IsConnected(ctx); err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	event := new(model.UserEvent)
	converter.UserToEvent(user, event)
	if err := u.UserProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	token, err := u.signAccessToken(ctx, user.ID)
	if err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Login", err)
	}

	user.Token = token

	res := new(model.UserResponse)
	converter.UserToTokenResponse(user, res)

	return res, nil
}
