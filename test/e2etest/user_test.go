package e2etest

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		Username: "khannedy",
		Password: "rahasia",
		Name:     "Eko Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, requestBody.Username, responseBody.Data.Username)
	assert.True(t, responseBody.Data.ID > 0)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestRegisterError(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		Username: "",
		Password: "",
		Name:     "",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestRegisterDuplicate(t *testing.T) {
	ClearAll()
	registerDefaultUser(t)

	requestBody := model.RegisterUserRequest{
		Username: "khannedy",
		Password: "rahasia",
		Name:     "Eko Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusConflict, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestLogin(t *testing.T) {
	ClearAll()
	registerDefaultUser(t)

	requestBody := model.LoginUserRequest{
		Username: "khannedy",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_login", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserLoginResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.NotEmpty(t, responseBody.Data.Token)

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(responseBody.Data.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viperConfig.GetString(configkey.AuthJWTSecret)), nil
	})
	assert.Nil(t, err)
	assert.True(t, token.Valid)
	user := new(entity.User)
	err = db.Where("username = ?", requestBody.Username).First(user).Error
	assert.Nil(t, err)
	assert.Equal(t, strconv.FormatInt(user.ID, 10), claims.Subject)
	assert.Equal(t, viperConfig.GetString(configkey.AuthJWTIssuer), claims.Issuer)
	assert.NotNil(t, claims.ExpiresAt)
}

func TestLoginWrongUsername(t *testing.T) {
	ClearAll()
	registerDefaultUser(t)

	requestBody := model.LoginUserRequest{
		Username: "wrong",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_login", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestLoginWrongPassword(t *testing.T) {
	ClearAll()
	registerDefaultUser(t)

	requestBody := model.LoginUserRequest{
		Username: "khannedy",
		Password: "wrong",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_login", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestGetCurrentUser(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)

	user := new(entity.User)
	err := db.Where("username = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/api/users/_current", nil)
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, user.Name, responseBody.Data.Name)
	assert.Equal(t, user.CreatedAt, responseBody.Data.CreatedAt)
	assert.Equal(t, user.UpdatedAt, responseBody.Data.UpdatedAt)
}

func TestGetCurrentUserFailed(t *testing.T) {
	ClearAll()

	registerDefaultUser(t)

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/api/users/_current", nil)
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "wrong")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestUpdateUserName(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)

	user := new(entity.User)
	err := db.Where("username = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	requestBody := model.UpdateUserRequest{
		Name: "Eko Kurniawan Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:3000/api/users/_current", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateUserPassword(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)

	user := new(entity.User)
	err := db.Where("username = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	requestBody := model.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:3000/api/users/_current", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.Equal(t, user.ID, responseBody.Data.ID)
	assert.NotNil(t, responseBody.Data.CreatedAt)
	assert.NotNil(t, responseBody.Data.UpdatedAt)

	user = new(entity.User)
	err = db.Where("username = ?", "khannedy").First(user).Error
	assert.Nil(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	assert.Nil(t, err)
}

func TestUpdateFailed(t *testing.T) {
	ClearAll()
	registerDefaultUser(t)

	requestBody := model.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:3000/api/users/_current", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "wrong")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)

	assert.NotNil(t, responseBody.ErrorMessage)
}

func TestFollowUser(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)

	// Register user to be followed
	followingUsername := "johndoe"
	registerUser(t, followingUsername, "rahasia", "John Doe")

	followingUser := new(entity.User)
	err := db.Where("username = ?", followingUsername).First(followingUser).Error
	assert.Nil(t, err)

	requestBody := model.FollowUserRequest{
		FollowingID: followingUser.ID,
	}
	bodyJson, err := json.Marshal(requestBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_follow", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	assert.Nil(t, err)
	assert.Equal(t, "ok", responseBody.Data)

	followerUser := new(entity.User)
	err = db.Where("username = ?", defaultUsername).First(followerUser).Error
	assert.Nil(t, err)

	follow := new(entity.Follow)
	err = db.Where("follower_id = ? AND following_id = ?", followerUser.ID, followingUser.ID).First(follow).Error
	assert.Nil(t, err)
}
