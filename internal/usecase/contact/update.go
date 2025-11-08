package contact

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
)

func (u *ContactUsecaseImpl) Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error) {
	if err := x.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Update", err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIDAndUserID(ctx, u.DB.WithContext(ctx), contact, req.ID, req.UserID); err != nil {
		return nil, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Update", err)
	}

	converter.ModelUpdateContactRequestToEntityContact(req, contact)

	if err := u.ContactRepository.Update(ctx, u.DB.WithContext(ctx), contact); err != nil {
		return nil, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Update", err)
	}

	if _, err := u.SlackClient.IsConnected(ctx, model.SlackIsConnectedRequest{}); err != nil {
		return nil, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Update", err)
	}

	event := new(model.ContactEvent)
	converter.EntityContactToModelContactEvent(contact, event)
	if err := u.ContactProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Update", err)
	}

	res := new(model.ContactResponse)
	converter.EntityContactToModelContactResponse(contact, res)

	return res, nil
}
