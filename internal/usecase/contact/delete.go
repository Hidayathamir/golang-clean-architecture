package contact

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *ContactUsecaseImpl) Delete(ctx context.Context, req *model.DeleteContactRequest) error {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, u.DB.WithContext(ctx), contact, req.ID, req.UserId); err != nil {
		return errkit.AddFuncName(err)
	}

	if err := u.ContactRepository.Delete(ctx, u.DB.WithContext(ctx), contact); err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
