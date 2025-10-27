package user2

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *User2UsecaseImpl) Profile(ctx context.Context, req *model.GetUser2Request) (*model.User2Response, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).Profile", err)
	}

	user := new(entity.User2)
	if err := u.User2Repository.FindByID(ctx, u.DB.WithContext(ctx), user, req.ID); err != nil {
		return nil, errkit.AddFuncName("user2.(*User2UsecaseImpl).Profile", err)
	}

	res := new(model.User2Response)
	converter.EntityUser2ToModelUser2Response(user, res)
	return res, nil
}
