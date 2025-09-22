package contact

import (
	"context"
	"golang-clean-architecture/internal/entity"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/model/converter"
	"golang-clean-architecture/pkg/errkit"
)

func (u *ContactUseCaseImpl) Update(ctx context.Context, req *model.UpdateContactRequest) (*model.ContactResponse, error) {
	tx := u.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, tx, contact, req.ID, req.UserId); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName(err)
	}

	contact.FirstName = req.FirstName
	contact.LastName = req.LastName
	contact.Email = req.Email
	contact.Phone = req.Phone

	if err := u.ContactRepository.Update(ctx, tx, contact); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errkit.AddFuncName(err)
	}

	if _, err := u.SlackClient.IsConnected(ctx); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	event := converter.ContactToEvent(contact)
	if err := u.ContactProducer.Send(ctx, event); err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return converter.ContactToResponse(contact), nil
}
