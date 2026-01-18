package e2etest

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/golang-jwt/jwt/v5"
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
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
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
	requestBody := model.RegisterUserRequest{
		Username: "",
		Password: "",
		Name:     "",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
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
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusConflict, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
}

func TestLogin(t *testing.T) {
	ClearAll()
	registerDefaultUser(t)

	requestBody := model.LoginUserRequest{
		Username: "khannedy",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_login", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserLoginResponse])
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotEmpty(t, responseBody.Data.Token)

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(responseBody.Data.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(viperConfig.GetString(configkey.AuthJWTSecret)), nil
	})
	require.Nil(t, err)
	require.True(t, token.Valid)
	user := new(entity.User)
	err = db.Where("username = ?", requestBody.Username).First(user).Error
	require.Nil(t, err)
	require.Equal(t, strconv.FormatInt(user.ID, 10), claims.Subject)
	require.Equal(t, viperConfig.GetString(configkey.AuthJWTIssuer), claims.Issuer)
	require.NotNil(t, claims.ExpiresAt)
}

func TestLoginWrongUsername(t *testing.T) {
	ClearAll()
	registerDefaultUser(t)

	requestBody := model.LoginUserRequest{
		Username: "wrong",
		Password: "rahasia",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_login", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
}

func TestLoginWrongPassword(t *testing.T) {
	ClearAll()
	registerDefaultUser(t)

	requestBody := model.LoginUserRequest{
		Username: "khannedy",
		Password: "wrong",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_login", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
}

func TestGetCurrentUser(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)

	user := new(entity.User)
	err := db.Where("username = ?", "khannedy").First(user).Error
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/api/users/_current", nil)
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.Equal(t, user.ID, responseBody.Data.ID)
	require.Equal(t, user.Name, responseBody.Data.Name)
	require.Equal(t, user.CreatedAt, responseBody.Data.CreatedAt)
	require.Equal(t, user.UpdatedAt, responseBody.Data.UpdatedAt)
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

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
}

func TestUpdateUserName(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)

	user := new(entity.User)
	err := db.Where("username = ?", "khannedy").First(user).Error
	require.Nil(t, err)

	requestBody := model.UpdateUserRequest{
		Name: "Eko Kurniawan Khannedy",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:3000/api/users/_current", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.Equal(t, user.ID, responseBody.Data.ID)
	require.Equal(t, requestBody.Name, responseBody.Data.Name)
	require.NotNil(t, responseBody.Data.CreatedAt)
	require.NotNil(t, responseBody.Data.UpdatedAt)
}

func TestUpdateUserPassword(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)

	user := new(entity.User)
	err := db.Where("username = ?", "khannedy").First(user).Error
	require.Nil(t, err)

	requestBody := model.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:3000/api/users/_current", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.Equal(t, user.ID, responseBody.Data.ID)
	require.NotNil(t, responseBody.Data.CreatedAt)
	require.NotNil(t, responseBody.Data.UpdatedAt)

	user = new(entity.User)
	err = db.Where("username = ?", "khannedy").First(user).Error
	require.Nil(t, err)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
	require.Nil(t, err)
}

func TestUpdateFailed(t *testing.T) {
	ClearAll()
	registerDefaultUser(t)

	requestBody := model.UpdateUserRequest{
		Password: "rahasialagi",
	}

	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:3000/api/users/_current", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "wrong")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[model.UserResponse])
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.NotNil(t, responseBody.ErrorMessage)
}

func TestFollowUser(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)

	// Register user to be followed
	followingUsername := "johndoe"
	registerUser(t, followingUsername, "rahasia", "John Doe")

	followingUser := new(entity.User)
	err := db.Where("username = ?", followingUsername).First(followingUser).Error
	require.Nil(t, err)

	requestBody := model.FollowUserRequest{
		FollowingID: followingUser.ID,
	}
	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_follow", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := new(response.WebResponse[string])
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)
	require.Equal(t, "ok", responseBody.Data)

	followerUser := new(entity.User)
	err = db.Where("username = ?", defaultUsername).First(followerUser).Error
	require.Nil(t, err)

	follow := new(entity.Follow)
	err = db.Where("follower_id = ? AND following_id = ?", followerUser.ID, followingUser.ID).First(follow).Error
	require.Nil(t, err)
}

func TestUserFollowScenario(t *testing.T) {
	ClearAll()

	// Register and login users
	tokenA := registerAndLoginUser(t, "user_a", "password", "User A")
	tokenB := registerAndLoginUser(t, "user_b", "password", "User B")
	tokenC := registerAndLoginUser(t, "user_c", "password", "User C")

	// Get User IDs
	userA := new(entity.User)
	err := db.Where("username = ?", "user_a").First(userA).Error
	require.Nil(t, err)

	userB := new(entity.User)
	err = db.Where("username = ?", "user_b").First(userB).Error
	require.Nil(t, err)

	userC := new(entity.User)
	err = db.Where("username = ?", "user_c").First(userC).Error
	require.Nil(t, err)

	// 1. User A follows User B
	followUser(t, tokenA, userB.ID)

	// 2. User C follows User B
	followUser(t, tokenC, userB.ID)

	// 3. User B follows User A
	followUser(t, tokenB, userA.ID)

	// Verify Database
	checkFollow(t, userA.ID, userB.ID)
	checkFollow(t, userC.ID, userB.ID)
	checkFollow(t, userB.ID, userA.ID)

	// need to wait for topic to be consume
	time.Sleep(5 * time.Second)

	err = db.Where("username = ?", "user_a").First(userA).Error
	require.Nil(t, err)
	require.Equal(t, 1, userA.FollowerCount)
	require.Equal(t, 1, userA.FollowingCount)

	err = db.Where("username = ?", "user_b").First(userB).Error
	require.Nil(t, err)
	require.Equal(t, 2, userB.FollowerCount)
	require.Equal(t, 1, userB.FollowingCount)

	err = db.Where("username = ?", "user_c").First(userC).Error
	require.Nil(t, err)
	require.Equal(t, 0, userC.FollowerCount)
	require.Equal(t, 1, userC.FollowingCount)
}
