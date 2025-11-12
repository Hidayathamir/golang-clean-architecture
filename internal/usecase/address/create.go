package address

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *AddressUsecaseImpl) Create(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIDAndUserID(ctx, u.DB.WithContext(ctx), contact, req.ContactID, req.UserID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	address := new(entity.Address)
	converter.ModelCreateAddressRequestToEntityAddress(req, address, contact.ID)

	if err := u.AddressRepository.Create(ctx, u.DB.WithContext(ctx), address); err != nil {
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
