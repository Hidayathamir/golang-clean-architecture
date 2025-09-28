package address

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/pkg/errkit"
)

func (u *AddressUsecaseImpl) List(ctx context.Context, req *model.ListAddressRequest) ([]model.AddressResponse, error) {
	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, u.DB.WithContext(ctx), contact, req.ContactId, req.UserId); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	addresses, err := u.AddressRepository.FindAllByContactId(ctx, u.DB.WithContext(ctx), contact.ID)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := make([]model.AddressResponse, len(addresses))
	for i, address := range addresses {
		res[i] = *converter.AddressToResponse(&address)
	}

	return res, nil
}
