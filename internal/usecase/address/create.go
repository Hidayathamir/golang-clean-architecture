package address

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/google/uuid"
)

func (u *AddressUsecaseImpl) Create(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("address.(*AddressUsecaseImpl).Create", err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, u.DB.WithContext(ctx), contact, req.ContactId, req.UserId); err != nil {
		return nil, errkit.AddFuncName("address.(*AddressUsecaseImpl).Create", err)
	}

	address := &entity.Address{
		ID:         uuid.NewString(),
		ContactId:  contact.ID,
		Street:     req.Street,
		City:       req.City,
		Province:   req.Province,
		PostalCode: req.PostalCode,
		Country:    req.Country,
	}

	if err := u.AddressRepository.Create(ctx, u.DB.WithContext(ctx), address); err != nil {
		return nil, errkit.AddFuncName("address.(*AddressUsecaseImpl).Create", err)
	}

	event := converter.AddressToEvent(address)
	if err := u.AddressProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName("address.(*AddressUsecaseImpl).Create", err)
	}

	return converter.AddressToResponse(address), nil
}
