package e2etest

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultUsername     = "khannedy"
	defaultUserPassword = "rahasia"
	defaultUserName     = "Eko Khannedy"
)

func registerDefaultUser(t *testing.T) {
	t.Helper()

	registerUser(t, defaultUsername, defaultUserPassword, defaultUserName)
}

func registerUser(t *testing.T, username, password, name string) {
	t.Helper()

	requestBody := model.RegisterUserRequest{
		Username: username,
		Password: password,
		Name:     name,
	}

	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users", strings.NewReader(string(bodyJSON)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)

	defer func() { _ = res.Body.Close() }()

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, requestBody.Username, responseBody.Data.Username)
}

func loginDefaultUser(t *testing.T) string {
	t.Helper()

	return loginUser(t, defaultUsername, defaultUserPassword)
}

func loginUser(t *testing.T, username, password string) string {
	t.Helper()

	requestBody := model.LoginUserRequest{
		Username: username,
		Password: password,
	}

	bodyJSON, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_login", strings.NewReader(string(bodyJSON)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)

	defer func() { _ = res.Body.Close() }()

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserLoginResponse])
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

func registerAndLoginUser(t *testing.T, username, password, name string) string {
	t.Helper()

	registerUser(t, username, password, name)
	return loginUser(t, username, password)
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

func ClearAll() {
	ClearUsers()
}

func ClearUsers() {
	err := db.Unscoped().Where("id is not null").Delete(&entity.User{}).Error
	if err != nil {
		x.Logger.Panicf("Failed clear user data : %+v", err)
	}
}

func GetFirstUser(t *testing.T) *entity.User {
	user := new(entity.User)
	err := db.First(user).Error
	assert.Nil(t, err)
	return user
}
