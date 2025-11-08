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
	if err := u.ContactRepository.FindByIDAndUserID(ctx, u.DB.WithContext(ctx), contact, req.ContactID, req.UserID); err != nil {
		return errkit.AddFuncName("address.(*AddressUsecaseImpl).Delete", err)
	}

	address := new(entity.Address)
	if err := u.AddressRepository.FindByIDAndContactID(ctx, u.DB.WithContext(ctx), address, req.ID, req.ContactID); err != nil {
		return errkit.AddFuncName("address.(*AddressUsecaseImpl).Delete", err)
	}

	err := u.DB.Transaction(func(tx *gorm.DB) error {
		if _, err := u.PaymentClient.Refund(ctx, model.PaymentRefundRequest{
			TransactionID: address.ID,
		}); err != nil {
			return errkit.AddFuncName("address.(*AddressUsecaseImpl).Delete", err)
		}

		if err := u.AddressRepository.Delete(ctx, tx, address); err != nil {
			return errkit.AddFuncName("address.(*AddressUsecaseImpl).Delete", err)
		}

		return nil
	})
	if err != nil {
		return errkit.AddFuncName("address.(*AddressUsecaseImpl).Delete", err)
	}

	return nil
}
