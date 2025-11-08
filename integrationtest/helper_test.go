package integrationtest

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	defaultUserID       = "khannedy"
	defaultUserPassword = "rahasia"
	defaultUserName     = "Eko Khannedy"
)

func registerDefaultUser(t *testing.T) {
	t.Helper()

	requestBody := model.RegisterUserRequest{
		ID:       defaultUserID,
		Password: defaultUserPassword,
		Name:     defaultUserName,
	}

	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/users", strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := app.Test(req)
	assert.Nil(t, err)

	defer func() { _ = res.Body.Close() }()

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, requestBody.ID, responseBody.Data.ID)
}

func loginDefaultUser(t *testing.T) string {
	t.Helper()

	requestBody := model.LoginUserRequest{
		ID:       defaultUserID,
		Password: defaultUserPassword,
	}

	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/users/_login", strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := app.Test(req)
	assert.Nil(t, err)

	defer func() { _ = res.Body.Close() }()

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.NotEmpty(t, responseBody.Data.Token)

	return responseBody.Data.Token
}

func registerAndLoginDefaultUser(t *testing.T) string {
	t.Helper()

	registerDefaultUser(t)
	return loginDefaultUser(t)
}

func bearerToken(token string) string {
	return "Bearer " + token
}

func loginAndGetDefaultUser(t *testing.T) (string, *entity.User) {
	t.Helper()

	token := registerAndLoginDefaultUser(t)
	user := GetFirstUser(t)
	return token, user
}

func createTodoViaAPI(t *testing.T, token string, requestBody model.CreateTodoRequest) model.TodoResponse {
	t.Helper()

	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	defer func() { _ = res.Body.Close() }()

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.TodoResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	return responseBody.Data
}

func createContactViaAPI(t *testing.T, token string, requestBody model.CreateContactRequest) model.ContactResponse {
	t.Helper()

	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/contacts", strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	defer func() { _ = res.Body.Close() }()

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.ContactResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	return responseBody.Data
}

func createAddressViaAPI(t *testing.T, token string, contactID string, requestBody model.CreateAddressRequest) model.AddressResponse {
	t.Helper()

	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/contacts/"+contactID+"/addresses", strings.NewReader(string(bodyJSON)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := app.Test(req)
	assert.Nil(t, err)

	defer func() { _ = res.Body.Close() }()

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.AddressResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	return responseBody.Data
}

func ClearAll() {
	ClearAddresses()
	ClearContact()
	ClearTodos()
	ClearUsers()
}

func ClearUsers() {
	err := db.Where("id is not null").Delete(&entity.User{}).Error
	if err != nil {
		x.Logger.Panicf("Failed clear user data : %+v", err)
	}
}

func ClearContact() {
	err := db.Where("id is not null").Delete(&entity.Contact{}).Error
	if err != nil {
		x.Logger.Panicf("Failed clear contact data : %+v", err)
	}
}

func ClearAddresses() {
	err := db.Where("id is not null").Delete(&entity.Address{}).Error
	if err != nil {
		x.Logger.Panicf("Failed clear address data : %+v", err)
	}
}

func ClearTodos() {
	err := db.Where("id is not null").Delete(&entity.Todo{}).Error
	if err != nil {
		x.Logger.Panicf("Failed clear todo data : %+v", err)
	}
}

func CreateContacts(user *entity.User, total int) {
	for i := range total {
		contact := &entity.Contact{
			ID:        uuid.NewString(),
			FirstName: "Contact",
			LastName:  strconv.Itoa(i),
			Email:     "contact" + strconv.Itoa(i) + "@example.com",
			Phone:     "08000000" + strconv.Itoa(i),
			UserID:    user.ID,
		}
		err := db.Create(contact).Error
		if err != nil {
			x.Logger.Panicf("Failed create contact data : %+v", err)
		}
	}
}

func CreateAddresses(t *testing.T, contact *entity.Contact, total int) {
	for range total {
		address := &entity.Address{
			ID:         uuid.NewString(),
			ContactID:  contact.ID,
			Street:     "Jalan Belum Jadi",
			City:       "Jakarta",
			Province:   "DKI Jakarta",
			PostalCode: "2131323",
			Country:    "Indonesia",
		}
		err := db.Create(address).Error
		assert.Nil(t, err)
	}
}

func GetFirstUser(t *testing.T) *entity.User {
	user := new(entity.User)
	err := db.First(user).Error
	assert.Nil(t, err)
	return user
}

func GetFirstContact(t *testing.T, user *entity.User) *entity.Contact {
	contact := new(entity.Contact)
	err := db.Where("user_id = ?", user.ID).First(contact).Error
	assert.Nil(t, err)
	return contact
}

func GetFirstAddress(t *testing.T, contact *entity.Contact) *entity.Address {
	address := new(entity.Address)
	err := db.Where("contact_id = ?", contact.ID).First(address).Error
	assert.Nil(t, err)
	return address
}

func CreateTodos(t *testing.T, user *entity.User, total int) {
	for i := range total {
		todo := &entity.Todo{
			ID:          uuid.NewString(),
			UserID:      user.ID,
			Title:       "Todo " + strconv.Itoa(i),
			Description: "Description " + strconv.Itoa(i),
			IsCompleted: false,
		}
		err := db.Create(todo).Error
		assert.Nil(t, err)
	}
}

func GetFirstTodo(t *testing.T, user *entity.User) *entity.Todo {
	todo := new(entity.Todo)
	err := db.Where("user_id = ?", user.ID).First(todo).Error
	assert.Nil(t, err)
	return todo
}
