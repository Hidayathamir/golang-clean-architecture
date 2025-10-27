package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/google/uuid"
)

func EntityAddressToModelAddressResponse(address *entity.Address, res *model.AddressResponse) {
	if address == nil || res == nil {
		return
	}

	res.ID = address.ID
	res.Street = address.Street
	res.City = address.City
	res.Province = address.Province
	res.PostalCode = address.PostalCode
	res.Country = address.Country
	res.CreatedAt = address.CreatedAt
	res.UpdatedAt = address.UpdatedAt
}

func EntityAddressToModelAddressEvent(address *entity.Address, event *model.AddressEvent) {
	if address == nil || event == nil {
		return
	}

	event.ID = address.ID
	event.ContactID = address.ContactID
	event.Street = address.Street
	event.City = address.City
	event.Province = address.Province
	event.PostalCode = address.PostalCode
	event.Country = address.Country
	event.CreatedAt = address.CreatedAt
	event.UpdatedAt = address.UpdatedAt
}

func ModelCreateAddressRequestToEntityAddress(req *model.CreateAddressRequest, address *entity.Address, contactID string) {
	if req == nil || address == nil {
		return
	}

	address.ID = uuid.NewString()
	address.ContactID = contactID
	address.Street = req.Street
	address.City = req.City
	address.Province = req.Province
	address.PostalCode = req.PostalCode
	address.Country = req.Country
}

func ModelUpdateAddressRequestToEntityAddress(req *model.UpdateAddressRequest, address *entity.Address) {
	address.Street = req.Street
	address.City = req.City
	address.Province = req.Province
	address.PostalCode = req.PostalCode
	address.Country = req.Country
}

func EntityAddressListToModelAddressResponseList(addresses entity.AddressList, res model.AddressResponseList) {
	if res == nil {
		return
	}

	for i := range addresses {
		EntityAddressToModelAddressResponse(&addresses[i], &res[i])
	}
}
