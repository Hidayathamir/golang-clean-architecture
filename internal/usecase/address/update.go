package address

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/pkg/errkit"
)

func (u *AddressUsecaseImpl) Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, u.DB.WithContext(ctx), contact, req.ContactId, req.UserId); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	address := new(entity.Address)
	if err := u.AddressRepository.FindByIdAndContactId(ctx, u.DB.WithContext(ctx), address, req.ID, contact.ID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	address.Street = req.Street
	address.City = req.City
	address.Province = req.Province
	address.PostalCode = req.PostalCode
	address.Country = req.Country

	if err := u.AddressRepository.Update(ctx, u.DB.WithContext(ctx), address); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.AddressToEvent(address)
	if err := u.AddressProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.AddressToResponse(address), nil
}
