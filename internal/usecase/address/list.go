package address

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/pkg/errkit"
)

func (u *AddressUsecaseImpl) List(ctx context.Context, req *model.ListAddressRequest) ([]model.AddressResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, tx, contact, req.ContactId, req.UserId); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	addresses, err := u.AddressRepository.FindAllByContactId(ctx, tx, contact.ID)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := make([]model.AddressResponse, len(addresses))
	for i, address := range addresses {
		res[i] = *converter.AddressToResponse(&address)
	}

	return res, nil
}
