package address

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/pkg/errkit"
)

func (u *AddressUsecaseImpl) Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error) {
	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, u.DB.WithContext(ctx), contact, req.ContactId, req.UserId); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	address := new(entity.Address)
	if err := u.AddressRepository.FindByIdAndContactId(ctx, u.DB.WithContext(ctx), address, req.ID, req.ContactId); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.AddressToResponse(address), nil
}
