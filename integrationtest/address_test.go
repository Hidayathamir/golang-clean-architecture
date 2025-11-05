package integrationtest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateAddress(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)

	requestBody := model.CreateAddressRequest{
		Street:     "Jalan Belum Jadi",
		City:       "Jakarta",
		Province:   "DKI Jakarta",
		PostalCode: "343443",
		Country:    "Indonesia",
	}

	resData := createAddressViaAPI(t, token, contact.ID, requestBody)

	assert.Equal(t, requestBody.Street, resData.Street)
	assert.Equal(t, requestBody.City, resData.City)
	assert.Equal(t, requestBody.Province, resData.Province)
	assert.Equal(t, requestBody.Country, resData.Country)
	assert.Equal(t, requestBody.PostalCode, resData.PostalCode)
	assert.NotEmpty(t, resData.ID)
	assert.NotZero(t, resData.CreatedAt)
	assert.NotZero(t, resData.UpdatedAt)
}

func TestCreateAddressFailed(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)

	requestBody := model.CreateAddressRequest{
		Street:     "Jalan Belum Jadi",
		City:       "Jakarta",
		Province:   "DKI Jakarta",
		PostalCode: "343443343443343443343443343443343443343443",
		Country:    "Indonesia",
	}
	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/contacts/"+contact.ID+"/addresses", strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestListAddresses(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)

	CreateAddresses(t, contact, 5)

	req := httptest.NewRequest(http.MethodGet, "/api/contacts/"+contact.ID+"/addresses", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponseList])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 5, len(responseBody.Data))
}

func TestListAddressesFailed(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)

	CreateAddresses(t, contact, 5)

	req := httptest.NewRequest(http.MethodGet, "/api/contacts/wrong/addresses", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponseList])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestGetAddress(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)
	CreateAddresses(t, contact, 1)
	address := GetFirstAddress(t, contact)

	req := httptest.NewRequest(http.MethodGet, "/api/contacts/"+contact.ID+"/addresses/"+address.ID, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, address.ID, responseBody.Data.ID)
	assert.Equal(t, address.Street, responseBody.Data.Street)
	assert.Equal(t, address.City, responseBody.Data.City)
	assert.Equal(t, address.Province, responseBody.Data.Province)
	assert.Equal(t, address.Country, responseBody.Data.Country)
	assert.Equal(t, address.PostalCode, responseBody.Data.PostalCode)
	assert.Equal(t, address.CreatedAt, responseBody.Data.CreatedAt)
	assert.Equal(t, address.UpdatedAt, responseBody.Data.UpdatedAt)
}

func TestGetAddressFailed(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)
	CreateAddresses(t, contact, 1)

	req := httptest.NewRequest(http.MethodGet, "/api/contacts/"+contact.ID+"/addresses/wrong", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestUpdateAddress(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)
	CreateAddresses(t, contact, 1)
	address := GetFirstAddress(t, contact)

	requestBody := model.UpdateAddressRequest{
		Street:     "Jalan Lagi Dijieun",
		City:       "Bandung",
		Province:   "Jawa Barat",
		PostalCode: "343443",
		Country:    "Indonesia",
	}
	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/contacts/"+contact.ID+"/addresses/"+address.ID, strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, requestBody.Street, responseBody.Data.Street)
	assert.Equal(t, requestBody.City, responseBody.Data.City)
	assert.Equal(t, requestBody.Province, responseBody.Data.Province)
	assert.Equal(t, requestBody.Country, responseBody.Data.Country)
	assert.Equal(t, requestBody.PostalCode, responseBody.Data.PostalCode)
	assert.NotEmpty(t, responseBody.Data.ID)
	assert.NotZero(t, responseBody.Data.CreatedAt)
	assert.NotZero(t, responseBody.Data.UpdatedAt)
}

func TestUpdateAddressFailed(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)
	CreateAddresses(t, contact, 1)
	address := GetFirstAddress(t, contact)

	requestBody := model.UpdateAddressRequest{
		Street:     "Jalan Lagi Dijieun",
		City:       "Bandung",
		Province:   "Jawa Barat",
		PostalCode: "343443343443343443343443343443343443343443",
		Country:    "Indonesia",
	}
	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/contacts/"+contact.ID+"/addresses/"+address.ID, strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestDeleteAddress(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)
	CreateAddresses(t, contact, 1)
	address := GetFirstAddress(t, contact)

	req := httptest.NewRequest(http.MethodDelete, "/api/contacts/"+contact.ID+"/addresses/"+address.ID, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, responseBody.Data)
}

func TestDeleteAddressFailed(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)
	CreateAddresses(t, contact, 1)

	req := httptest.NewRequest(http.MethodDelete, "/api/contacts/"+contact.ID+"/addresses/wrong", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}
