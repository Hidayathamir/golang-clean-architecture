package address

import (
	"context"
	"strconv"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"gorm.io/gorm"
)

func (u *AddressUsecaseImpl) Delete(ctx context.Context, req *model.DeleteAddressRequest) error {
	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIDAndUserID(ctx, u.DB.WithContext(ctx), contact, req.ContactID, req.UserID); err != nil {
		return errkit.AddFuncName(err)
	}

	address := new(entity.Address)
	if err := u.AddressRepository.FindByIDAndContactID(ctx, u.DB.WithContext(ctx), address, req.ID, req.ContactID); err != nil {
		return errkit.AddFuncName(err)
	}

	err := u.DB.Transaction(func(tx *gorm.DB) error {
		if _, err := u.PaymentClient.Refund(ctx, model.PaymentRefundRequest{
			TransactionID: strconv.FormatInt(address.ID, 10),
		}); err != nil {
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
