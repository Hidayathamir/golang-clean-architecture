package integrationtest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/google/uuid"
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

func TestCreateAddressUnauthorized(t *testing.T) {
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
	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/contacts/"+contact.ID+"/addresses", strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestCreateAddressOtherUserContact(t *testing.T) {
	ClearAll()
	tokenUserA, userA := loginAndGetDefaultUser(t)

	createContactViaAPI(t, tokenUserA, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, userA)

	otherUsername := uuid.NewString()
	otherPassword := "secret"
	otherToken := registerAndLoginUser(t, otherUsername, otherPassword, "Other User")

	requestBody := model.CreateAddressRequest{
		Street:     "Hack Street",
		City:       "Hack City",
		Province:   "Hack State",
		PostalCode: "12345",
		Country:    "Nowhere",
	}
	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/contacts/"+contact.ID+"/addresses", strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(otherToken))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestListAddressesUnauthorized(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)

	req := httptest.NewRequest(http.MethodGet, "/api/contacts/"+contact.ID+"/addresses", nil)
	req.Header.Set("Accept", "application/json")

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponseList])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
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

func TestCreateAddressInvalidContactID(t *testing.T) {
	ClearAll()
	token, _ := loginAndGetDefaultUser(t)

	requestBody := model.CreateAddressRequest{
		Street:     "Jalan Belum Jadi",
		City:       "Jakarta",
		Province:   "DKI Jakarta",
		PostalCode: "343443",
		Country:    "Indonesia",
	}
	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/contacts/not-a-uuid/addresses", strings.NewReader(string(bodyJSON)))
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
	assert.NotNil(t, responseBody.ErrorMessage)
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

func TestGetAddressInvalidUUID(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)

	req := httptest.NewRequest(http.MethodGet, "/api/contacts/"+contact.ID+"/addresses/not-a-uuid", nil)
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
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestListAddressesOtherUser(t *testing.T) {
	ClearAll()
	tokenUserA, userA := loginAndGetDefaultUser(t)

	createContactViaAPI(t, tokenUserA, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, userA)
	CreateAddresses(t, contact, 1)

	otherUsername := uuid.NewString()
	otherPassword := "secret"
	otherToken := registerAndLoginUser(t, otherUsername, otherPassword, "Other User")

	req := httptest.NewRequest(http.MethodGet, "/api/contacts/"+contact.ID+"/addresses", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(otherToken))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponseList])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestGetAddressMismatchedContact(t *testing.T) {
	ClearAll()
	token, _ := loginAndGetDefaultUser(t)

	first := createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Owner",
		LastName:  "One",
		Email:     "owner1@example.com",
		Phone:     "08000000001",
	})
	second := createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Owner",
		LastName:  "Two",
		Email:     "owner2@example.com",
		Phone:     "08000000002",
	})
	CreateAddresses(t, &entity.Contact{ID: first.ID}, 1)
	address := GetFirstAddress(t, &entity.Contact{ID: first.ID})

	req := httptest.NewRequest(http.MethodGet, "/api/contacts/"+second.ID+"/addresses/"+address.ID, nil)
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
	assert.NotNil(t, responseBody.ErrorMessage)
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

func TestUpdateAddressUnauthorized(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)
	createAddress := createAddressViaAPI(t, token, contact.ID, model.CreateAddressRequest{})

	requestBody := model.UpdateAddressRequest{
		Street:     "Unauthorized Street",
		City:       "Unauthorized City",
		Province:   "Unauthorized Province",
		PostalCode: "11111",
		Country:    "Unauthorized",
	}
	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/contacts/"+contact.ID+"/addresses/"+createAddress.ID, strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestUpdateAddressOtherUser(t *testing.T) {
	ClearAll()
	tokenUserA, userA := loginAndGetDefaultUser(t)

	createContactViaAPI(t, tokenUserA, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, userA)
	address := createAddressViaAPI(t, tokenUserA, contact.ID, model.CreateAddressRequest{})

	otherToken := registerAndLoginUser(t, uuid.NewString(), "secret", "Other User")

	requestBody := model.UpdateAddressRequest{
		Street:     "Hack Street",
		City:       "Hack City",
		Province:   "Hack Province",
		PostalCode: "00000",
		Country:    "Hackland",
	}
	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/contacts/"+contact.ID+"/addresses/"+address.ID, strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(otherToken))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
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

func TestDeleteAddressUnauthorized(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)
	address := createAddressViaAPI(t, token, contact.ID, model.CreateAddressRequest{})

	req := httptest.NewRequest(http.MethodDelete, "/api/contacts/"+contact.ID+"/addresses/"+address.ID, nil)
	req.Header.Set("Accept", "application/json")

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestDeleteAddressInvalidUUID(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)

	req := httptest.NewRequest(http.MethodDelete, "/api/contacts/"+contact.ID+"/addresses/not-a-uuid", nil)
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
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestDeleteAddressOtherUser(t *testing.T) {
	ClearAll()
	tokenUserA, userA := loginAndGetDefaultUser(t)

	createContactViaAPI(t, tokenUserA, model.CreateContactRequest{
		FirstName: "Eko",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, userA)
	address := createAddressViaAPI(t, tokenUserA, contact.ID, model.CreateAddressRequest{})

	otherToken := registerAndLoginUser(t, uuid.NewString(), "secret", "Other User")

	req := httptest.NewRequest(http.MethodDelete, "/api/contacts/"+contact.ID+"/addresses/"+address.ID, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(otherToken))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[bool])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
}
