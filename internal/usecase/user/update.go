package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserUsecaseImpl) Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Update", err)
	}

	user := new(entity.User)
	if err := u.UserRepository.FindByID(ctx, u.DB.WithContext(ctx), user, req.ID); err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Update", err)
	}

	var password string
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Update", err)
		}
		password = string(hashedPassword)
	}

	converter.ModelUpdateUserRequestToEntityUser(req, user, password)

	if err := u.UserRepository.Update(ctx, u.DB.WithContext(ctx), user); err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Update", err)
	}

	event := new(model.UserEvent)
	converter.UserToEvent(user, event)
	if err := u.UserProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName("user.(*UserUsecaseImpl).Update", err)
	}

	res := new(model.UserResponse)
	converter.UserToResponse(user, res)

	return res, nil
}
