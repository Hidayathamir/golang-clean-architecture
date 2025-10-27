package address

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *AddressUsecaseImpl) Get(ctx context.Context, req *model.GetAddressRequest) (*model.AddressResponse, error) {
	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIDAndUserID(ctx, u.DB.WithContext(ctx), contact, req.ContactID, req.UserID); err != nil {
		return nil, errkit.AddFuncName("address.(*AddressUsecaseImpl).Get", err)
	}

	address := new(entity.Address)
	if err := u.AddressRepository.FindByIDAndContactID(ctx, u.DB.WithContext(ctx), address, req.ID, req.ContactID); err != nil {
		return nil, errkit.AddFuncName("address.(*AddressUsecaseImpl).Get", err)
	}

	res := new(model.AddressResponse)
	converter.EntityAddressToModelAddressResponse(address, res)

	return res, nil
}
