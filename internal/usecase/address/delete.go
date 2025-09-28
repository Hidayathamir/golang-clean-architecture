package address

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
)

func (u *AddressUsecaseImpl) Delete(ctx context.Context, req *model.DeleteAddressRequest) error {
	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, u.DB.WithContext(ctx), contact, req.ContactId, req.UserId); err != nil {
		return errkit.AddFuncName(err)
	}

	address := new(entity.Address)
	if err := u.AddressRepository.FindByIdAndContactId(ctx, u.DB.WithContext(ctx), address, req.ID, req.ContactId); err != nil {
		return errkit.AddFuncName(err)
	}

	err := u.DB.Transaction(func(tx *gorm.DB) error {
		if _, err := u.PaymentClient.Refund(ctx, address.ID); err != nil {
			return errkit.AddFuncName(err)
		}

		if err := u.AddressRepository.Delete(ctx, tx, address); err != nil {
			return errkit.AddFuncName(err)
		}

		return nil
	})
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}
