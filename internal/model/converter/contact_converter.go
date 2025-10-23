package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
)

func ContactToResponse(contact *entity.Contact) *model.ContactResponse {
	return &model.ContactResponse{
		ID:        contact.ID,
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}
}

func ContactsToResponses(contacts []entity.Contact) []model.ContactResponse {
	responses := make([]model.ContactResponse, 0, len(contacts))
	for _, contact := range contacts {
		responses = append(responses, *ContactToResponse(&contact))
	}
	return responses
}

func ContactToEvent(contact *entity.Contact) *model.ContactEvent {
	return &model.ContactEvent{
		ID:        contact.ID,
		UserID:    contact.UserId,
		FirstName: contact.FirstName,
		LastName:  contact.LastName,
		Email:     contact.Email,
		Phone:     contact.Phone,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
	}
}
