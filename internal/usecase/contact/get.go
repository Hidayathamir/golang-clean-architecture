package contact

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *ContactUsecaseImpl) Get(ctx context.Context, req *model.GetContactRequest) (*model.ContactResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Get", err)
	}

	contact := new(entity.Contact)
	if err := u.ContactRepository.FindByIdAndUserId(ctx, u.DB.WithContext(ctx), contact, req.ID, req.UserId); err != nil {
		return nil, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Get", err)
	}

	return converter.ContactToResponse(contact), nil
}
