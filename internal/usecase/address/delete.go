package address

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/pkg/errkit"
)

func (u *AddressUseCaseImpl) Delete(ctx context.Context, req *model.DeleteAddressRequest) error {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, tx, contact, req.ContactId, req.UserId); err != nil {
		return errkit.AddFuncName(err)
	}

	address := new(entity.Address)
	if err := u.AddressRepository.FindByIdAndContactId(ctx, tx, address, req.ID, req.ContactId); err != nil {
		return errkit.AddFuncName(err)
	}

	if _, err := u.PaymentClient.Refund(ctx, address.ID); err != nil {
		return errkit.AddFuncName(err)
	}

	if err := u.AddressRepository.Delete(ctx, tx, address); err != nil {
		return errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
