package user2

import (
	"context"
	"errors"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (u *User2UsecaseImpl) Register(ctx context.Context, req *model.RegisterUser2Request) (*model.User2Response, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).Register", err)
	}

	candidate := new(entity.User2)
	if err := u.User2Repository.FindByEmail(ctx, u.DB.WithContext(ctx), candidate, req.Email); err == nil {
		err = errors.New("email already registered")
		err = errkit.Conflict(err)
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).Register", err)
	} else {
		var httpErr *errkit.HTTPError
		if !errors.As(err, &httpErr) || httpErr.HTTPCode != http.StatusNotFound {
			return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).Register", err)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).Register", err)
	}

	user := new(entity.User2)
	converter.ModelRegisterUser2RequestToEntityUser2(req, user, string(hashedPassword), uuid.NewString())

	if err := u.User2Repository.Create(ctx, u.DB.WithContext(ctx), user); err != nil {
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).Register", err)
	}

	res := new(model.User2Response)
	converter.EntityUser2ToModelUser2Response(user, res)

	return res, nil
}
