package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/google/uuid"
)

func ModelCreateContactRequestToEntityContact(req *model.CreateContactRequest, contact *entity.Contact) {
	if req == nil || contact == nil {
		return
	}

	contact.ID = uuid.NewString()
	contact.FirstName = req.FirstName
	contact.LastName = req.LastName
	contact.Email = req.Email
	contact.Phone = req.Phone
	contact.UserID = req.UserID
}

func ModelUpdateContactRequestToEntityContact(req *model.UpdateContactRequest, contact *entity.Contact) {
	if req == nil || contact == nil {
		return
	}

	contact.FirstName = req.FirstName
	contact.LastName = req.LastName
	contact.Email = req.Email
	contact.Phone = req.Phone
}

func EntityContactToModelContactResponse(contact *entity.Contact, res *model.ContactResponse) {
	if contact == nil || res == nil {
		return
	}

	res.ID = contact.ID
	res.FirstName = contact.FirstName
	res.LastName = contact.LastName
	res.Email = contact.Email
	res.Phone = contact.Phone
	res.CreatedAt = contact.CreatedAt
	res.UpdatedAt = contact.UpdatedAt
}

func EntityContactListToModelContactResponseList(contacts entity.ContactList, res model.ContactResponseList) {
	if res == nil {
		return
	}

	for i := range contacts {
		EntityContactToModelContactResponse(&contacts[i], &res[i])
	}
}

func EntityContactToModelContactEvent(contact *entity.Contact, event *model.ContactEvent) {
	if contact == nil || event == nil {
		return
	}

	event.ID = contact.ID
	event.UserID = contact.UserID
	event.FirstName = contact.FirstName
	event.LastName = contact.LastName
	event.Email = contact.Email
	event.Phone = contact.Phone
	event.CreatedAt = contact.CreatedAt
	event.UpdatedAt = contact.UpdatedAt
}
