package address

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *AddressUsecaseImpl) Update(ctx context.Context, req *model.UpdateAddressRequest) (*model.AddressResponse, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIDAndUserID(ctx, u.DB.WithContext(ctx), contact, req.ContactID, req.UserID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	address := new(entity.Address)
	if err := u.AddressRepository.FindByIDAndContactID(ctx, u.DB.WithContext(ctx), address, req.ID, contact.ID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	converter.ModelUpdateAddressRequestToEntityAddress(req, address)

	if err := u.AddressRepository.Update(ctx, u.DB.WithContext(ctx), address); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := new(model.AddressEvent)
	converter.EntityAddressToModelAddressEvent(address, event)
	if err := u.AddressProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := new(model.AddressResponse)
	converter.EntityAddressToModelAddressResponse(address, res)

	return res, nil
}
