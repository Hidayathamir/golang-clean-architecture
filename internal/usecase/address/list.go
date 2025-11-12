package address

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *AddressUsecaseImpl) List(ctx context.Context, req *model.ListAddressRequest) (model.AddressResponseList, error) {
	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIDAndUserID(ctx, u.DB.WithContext(ctx), contact, req.ContactID, req.UserID); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	addresses, err := u.AddressRepository.FindAllByContactID(ctx, u.DB.WithContext(ctx), contact.ID)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	res := make(model.AddressResponseList, len(addresses))
	converter.EntityAddressListToModelAddressResponseList(addresses, res)

	return res, nil
}
