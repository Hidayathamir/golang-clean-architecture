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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateContact(t *testing.T) {
	ClearAll()
	token, _ := loginAndGetDefaultUser(t)

	requestBody := model.CreateContactRequest{
		FirstName: "Eko Kurniawan",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	}

	resData := createContactViaAPI(t, token, requestBody)

	assert.Equal(t, requestBody.FirstName, resData.FirstName)
	assert.Equal(t, requestBody.LastName, resData.LastName)
	assert.Equal(t, requestBody.Email, resData.Email)
	assert.Equal(t, requestBody.Phone, resData.Phone)
	assert.NotEmpty(t, resData.ID)
	assert.NotZero(t, resData.CreatedAt)
	assert.NotZero(t, resData.UpdatedAt)
}

func TestCreateContactFailed(t *testing.T) {
	ClearAll()
	token, _ := loginAndGetDefaultUser(t)

	requestBody := model.CreateContactRequest{}
	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/contacts", strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestGetConnect(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko Kurniawan",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)

	req := httptest.NewRequest(http.MethodGet, "/api/contacts/"+contact.ID, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, contact.ID, responseBody.Data.ID)
	assert.Equal(t, contact.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, contact.LastName, responseBody.Data.LastName)
	assert.Equal(t, contact.Email, responseBody.Data.Email)
	assert.Equal(t, contact.Phone, responseBody.Data.Phone)
	assert.Equal(t, contact.CreatedAt, responseBody.Data.CreatedAt)
	assert.Equal(t, contact.UpdatedAt, responseBody.Data.UpdatedAt)
}

func TestGetContactFailed(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko Kurniawan",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	_ = GetFirstContact(t, user)

	req := httptest.NewRequest(http.MethodGet, "/api/contacts/"+uuid.NewString(), nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestUpdateContact(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko Kurniawan",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)

	requestBody := model.UpdateContactRequest{
		FirstName: "Eko",
		LastName:  "Budiman",
		Email:     "budiman@example.com",
		Phone:     "089898989",
	}
	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/contacts/"+contact.ID, strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, requestBody.FirstName, responseBody.Data.FirstName)
	assert.Equal(t, requestBody.LastName, responseBody.Data.LastName)
	assert.Equal(t, requestBody.Email, responseBody.Data.Email)
	assert.Equal(t, requestBody.Phone, responseBody.Data.Phone)
	assert.NotEmpty(t, responseBody.Data.ID)
	assert.NotZero(t, responseBody.Data.CreatedAt)
	assert.NotZero(t, responseBody.Data.UpdatedAt)
}

func TestUpdateContactFailed(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko Kurniawan",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)

	requestBody := model.UpdateContactRequest{}
	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/contacts/"+contact.ID, strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
}

func TestDeleteContact(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko Kurniawan",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	contact := GetFirstContact(t, user)

	req := httptest.NewRequest(http.MethodDelete, "/api/contacts/"+contact.ID, nil)
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

func TestDeleteContactFailed(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	createContactViaAPI(t, token, model.CreateContactRequest{
		FirstName: "Eko Kurniawan",
		LastName:  "Khannedy",
		Email:     "eko@example.com",
		Phone:     "088888888888",
	})
	_ = GetFirstContact(t, user)

	req := httptest.NewRequest(http.MethodDelete, "/api/contacts/"+uuid.NewString(), nil)
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

func TestSearchContact(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	CreateContacts(user, 20)

	req := httptest.NewRequest(http.MethodGet, "/api/contacts", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.ContactResponseList])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 10, len(responseBody.Data))
	assert.Equal(t, int64(20), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(2), responseBody.Paging.TotalPage)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.Size)
}

func TestSearchContactWithPagination(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	CreateContacts(user, 20)

	req := httptest.NewRequest(http.MethodGet, "/api/contacts?page=2&size=5", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.ContactResponseList])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 5, len(responseBody.Data))
	assert.Equal(t, int64(20), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(4), responseBody.Paging.TotalPage)
	assert.Equal(t, 2, responseBody.Paging.Page)
	assert.Equal(t, 5, responseBody.Paging.Size)
}

func TestSearchContactWithFilter(t *testing.T) {
	ClearAll()
	token, user := loginAndGetDefaultUser(t)

	CreateContacts(user, 20)

	req := httptest.NewRequest(http.MethodGet, "/api/contacts?name=contact&phone=08000000&email=example.com", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.ContactResponseList])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, 10, len(responseBody.Data))
	assert.Equal(t, int64(20), responseBody.Paging.TotalItem)
	assert.Equal(t, int64(2), responseBody.Paging.TotalPage)
	assert.Equal(t, 1, responseBody.Paging.Page)
	assert.Equal(t, 10, responseBody.Paging.Size)
}
