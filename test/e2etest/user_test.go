package e2etest

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	ClearAll()

	// Prepare registration request
	requestBody := dto.RegisterUserRequest{
		Username: "khannedy",
		Password: "rahasia",
		Name:     "Eko Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Send registration request
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify response status code
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Verify response body
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.Equal(t, requestBody.Username, responseBody.Data.Username)
	require.True(t, responseBody.Data.ID > 0)
	require.Equal(t, requestBody.Name, responseBody.Data.Name)
	require.NotNil(t, responseBody.Data.CreatedAt)
	require.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestRegisterError(t *testing.T) {
	ClearAll()

	// Prepare invalid registration request (empty fields)
	requestBody := dto.RegisterUserRequest{
		Username: "",
		Password: "",
		Name:     "",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Send registration request
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code is Bad Request
	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	// Verify error message exists in response
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
}

func TestRegisterDuplicate(t *testing.T) {
	ClearAll()

	// Register user first
	registerDefaultUser(t)

	// Try to register the same user again
	requestBody := dto.RegisterUserRequest{
		Username: "khannedy",
		Password: "rahasia",
		Name:     "Eko Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code is Conflict
	require.Equal(t, http.StatusConflict, res.StatusCode)

	// Verify error message exists in response
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
}

func TestLogin(t *testing.T) {
	ClearAll()

	// Register user first
	registerDefaultUser(t)

	// Prepare login request
	requestBody := dto.LoginUserRequest{
		Username: "khannedy",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Send login request
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_login", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code is OK
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Verify response body contains token
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserLoginResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotEmpty(t, responseBody.Data.Token)

	// Verify JWT token claims
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(responseBody.Data.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.GetAuthJWTSecret()), nil
	})
	require.Nil(t, err)
	require.True(t, token.Valid)

	// Verify user data matches token subject
	user := &entity.User{}
	err = db.Where("username = ?", requestBody.Username).First(user).Error
	require.Nil(t, err)
	require.Equal(t, strconv.FormatInt(user.ID, 10), claims.Subject)
	require.Equal(t, cfg.GetAuthJWTIssuer(), claims.Issuer)
	require.NotNil(t, claims.ExpiresAt)
}

func TestLoginWrongUsername(t *testing.T) {
	ClearAll()

	// Register user first
	registerDefaultUser(t)

	// Prepare login request with wrong username
	requestBody := dto.LoginUserRequest{
		Username: "wrong",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Send login request
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_login", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code is Unauthorized
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// Verify error message exists
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
}

func TestLoginWrongPassword(t *testing.T) {
	ClearAll()

	// Register user first
	registerDefaultUser(t)

	// Prepare login request with wrong password
	requestBody := dto.LoginUserRequest{
		Username: "khannedy",
		Password: "wrong",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Send login request
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_login", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code is Unauthorized
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// Verify error message exists
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
}

func TestGetCurrentUser(t *testing.T) {
	ClearAll()

	// Register and login user
	token := registerAndLoginDefaultUser(t)

	// Get user data from database for comparison
	user := &entity.User{}
	err := db.Where("username = ?", "khannedy").First(user).Error
	require.Nil(t, err)

	// Send get current user request
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/api/users/_current", nil)
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code is OK
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Verify response body matches database record
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.Equal(t, user.ID, responseBody.Data.ID)
	require.Equal(t, user.Name, responseBody.Data.Name)
	require.Equal(t, user.CreatedAt, responseBody.Data.CreatedAt)
	require.Equal(t, user.UpdatedAt, responseBody.Data.UpdatedAt)
}

func TestGetCurrentUserFailed(t *testing.T) {
	ClearAll()

	// Register user
	registerDefaultUser(t)

	// Send get current user request with invalid authorization
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/api/users/_current", nil)
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "wrong")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code is Unauthorized
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// Verify error message exists
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
}

func TestUpdateUserName(t *testing.T) {
	ClearAll()

	// Register and login user
	token := registerAndLoginDefaultUser(t)

	// Get user data from database
	user := &entity.User{}
	err := db.Where("username = ?", "khannedy").First(user).Error
	require.Nil(t, err)

	// Prepare update request (name change)
	requestBody := dto.UpdateUserRequest{
		Name: "Eko Kurniawan Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Send update request
	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:3000/api/users/_current", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code is OK
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Verify response body matches update request
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.Equal(t, user.ID, responseBody.Data.ID)
	require.Equal(t, requestBody.Name, responseBody.Data.Name)
	require.NotNil(t, responseBody.Data.CreatedAt)
	require.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateUserPassword(t *testing.T) {
	ClearAll()

	// Register and login user
	token := registerAndLoginDefaultUser(t)

	// Get user data from database
	user := &entity.User{}
	err := db.Where("username = ?", "khannedy").First(user).Error
	require.Nil(t, err)

	// Prepare update request (password change)
	requestBody := dto.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Send update request
	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:3000/api/users/_current", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code is OK
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Verify response body
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.Equal(t, user.ID, responseBody.Data.ID)
	require.NotNil(t, responseBody.Data.CreatedAt)
	require.NotNil(t, responseBody.Data.UpdatedAt)

	// Verify new password is correctly hashed and stored in database
	user = &entity.User{}
	err = db.Where("username = ?", "khannedy").First(user).Error
	require.Nil(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	require.Nil(t, err)
}

func TestUpdateFailed(t *testing.T) {
	ClearAll()

	// Register user
	registerDefaultUser(t)

	// Prepare update request
	requestBody := dto.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Send update request with invalid authorization
	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:3000/api/users/_current", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "wrong")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code is Unauthorized
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	// Verify error message exists
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
}

func TestFollowUser(t *testing.T) {
	ClearAll()

	// Register and login follower user
	token := registerAndLoginDefaultUser(t)

	// Register user to be followed
	followingUsername := "johndoe"
	registerUser(t, followingUsername, "rahasia", "John Doe")

	followingUser := &entity.User{}
	err := db.Where("username = ?", followingUsername).First(followingUser).Error
	require.Nil(t, err)

	// Prepare follow request
	requestBody := dto.FollowUserRequest{
		FollowingID: followingUser.ID,
	}
	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Send follow request
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_follow", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code is OK
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Verify response body
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[string]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)
	require.Equal(t, "ok", responseBody.Data)

	// Verify follow relationship in database
	followerUser := &entity.User{}
	err = db.Where("username = ?", defaultUsername).First(followerUser).Error
	require.Nil(t, err)

	follow := &entity.Follow{}
	err = db.Where("follower_id = ? AND following_id = ?", followerUser.ID, followingUser.ID).First(follow).Error
	require.Nil(t, err)
}

func TestUserFollowScenario(t *testing.T) {
	ClearAll()

	// Register and login multiple users
	tokenA := registerAndLoginUser(t, "user_a", "password", "User A")
	tokenB := registerAndLoginUser(t, "user_b", "password", "User B")
	tokenC := registerAndLoginUser(t, "user_c", "password", "User C")

	// Get User IDs
	userA := &entity.User{}
	err := db.Where("username = ?", "user_a").First(userA).Error
	require.Nil(t, err)

	userB := &entity.User{}
	err = db.Where("username = ?", "user_b").First(userB).Error
	require.Nil(t, err)

	userC := &entity.User{}
	err = db.Where("username = ?", "user_c").First(userC).Error
	require.Nil(t, err)

	// 1. User A follows User B
	followUser(t, tokenA, userB.ID)

	// 2. User C follows User B
	followUser(t, tokenC, userB.ID)

	// 3. User B follows User A
	followUser(t, tokenB, userA.ID)

	// Verify follow relationships in Database
	checkFollow(t, userA.ID, userB.ID)
	checkFollow(t, userC.ID, userB.ID)
	checkFollow(t, userB.ID, userA.ID)
}
