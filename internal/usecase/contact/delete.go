package contact

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/errkit"
)

func (u *ContactUsecaseImpl) Delete(ctx context.Context, req *model.DeleteContactRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, tx, contact, req.ID, req.UserId); err != nil {
		return errkit.AddFuncName(err)
	}

	if err := u.ContactRepository.Delete(ctx, tx, contact); err != nil {
		return errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
