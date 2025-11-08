package contact

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ContactUsecaseImpl) Create(ctx context.Context, req *model.CreateContactRequest) (*model.ContactResponse, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Create", err)
	}

	contact := new(entity.Contact)
	converter.ModelCreateContactRequestToEntityContact(req, contact)

	if err := u.ContactRepository.Create(ctx, u.DB.WithContext(ctx), contact); err != nil {
		return nil, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Create", err)
	}

	event := new(model.ContactEvent)
	converter.EntityContactToModelContactEvent(contact, event)
	if err := u.ContactProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Create", err)
	}

	res := new(model.ContactResponse)
	converter.EntityContactToModelContactResponse(contact, res)

	return res, nil
}
