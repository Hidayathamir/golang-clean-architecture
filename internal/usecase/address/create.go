package address

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/pkg/errkit"

	"github.com/google/uuid"
)

func (u *AddressUsecaseImpl) Create(ctx context.Context, req *model.CreateAddressRequest) (*model.AddressResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, tx, contact, req.ContactId, req.UserId); err != nil {
		return nil, errkit.AddFuncName(err)
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

	if err := u.AddressRepository.Create(ctx, tx, address); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.AddressToEvent(address)
	if err := u.AddressProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.AddressToResponse(address), nil
}
