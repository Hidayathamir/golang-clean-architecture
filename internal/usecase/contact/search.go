package contact

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model/converter"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

func (u *ContactUsecaseImpl) Search(ctx context.Context, req *model.SearchContactRequest) (model.ContactResponseList, int64, error) {
	if err := u.Validate.Struct(req); err != nil {
		err = errkit.BadRequest(err)
		return nil, 0, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Search", err)
	}

	contacts, total, err := u.ContactRepository.Search(ctx, u.DB.WithContext(ctx), req)
	if err != nil {
		return nil, 0, errkit.AddFuncName("contact.(*ContactUsecaseImpl).Search", err)
	}

	res := make(model.ContactResponseList, len(contacts))
	converter.EntityContactListToModelContactResponseList(contacts, res)

	return res, total, nil
}
