package e2etest

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/x"
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

// registerUser performs an HTTP request to register a new user.
func registerUser(t *testing.T, username, password, name string) {
	t.Helper()

	// Prepare registration request body
	requestBody := dto.RegisterUserRequest{
		Username: username,
		Password: password,
		Name:     name,
	}

	bodyJSON, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Create and send POST /api/users request
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users", strings.NewReader(string(bodyJSON)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)

	defer func() { _ = res.Body.Close() }()

	// Verify status code and response data
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, requestBody.Username, responseBody.Data.Username)
}

func loginDefaultUser(t *testing.T) string {
	t.Helper()

	return loginUser(t, defaultUsername, defaultUserPassword)
}

// loginUser performs an HTTP request to login a user and returns the JWT token.
func loginUser(t *testing.T, username, password string) string {
	t.Helper()

	// Prepare login request body
	requestBody := dto.LoginUserRequest{
		Username: username,
		Password: password,
	}

	bodyJSON, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Create and send POST /api/users/_login request
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_login", strings.NewReader(string(bodyJSON)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)

	defer func() { _ = res.Body.Close() }()

	// Verify status code and retrieve token
	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UserLoginResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NotEmpty(t, responseBody.Data.Token)

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

// will remove when loginAndGetDefaultUser is actually used,
// this just so it not throw warning
var _ = loginAndGetDefaultUser

// followUser performs an HTTP request to follow another user.
func followUser(t *testing.T, token string, followingID int64) {
	t.Helper()

	// Prepare follow request body
	requestBody := dto.FollowUserRequest{
		FollowingID: followingID,
	}
	bodyJson, err := json.Marshal(requestBody)
	require.Nil(t, err)

	// Create and send POST /api/users/_follow request with Authorization header
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_follow", strings.NewReader(string(bodyJson)))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify status code and response data
	require.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[string]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)
	require.Equal(t, "ok", responseBody.Data)
}

// checkFollow verifies the existence of a follow relationship in the database.
func checkFollow(t *testing.T, followerID, followingID int64) {
	t.Helper()

	follow := &entity.Follow{}
	err := db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(follow).Error
	require.Nil(t, err)
	require.Equal(t, followerID, follow.FollowerID)
	require.Equal(t, followingID, follow.FollowingID)
}

func commentImage(t *testing.T, token string, imageID int64, comment string) {
	t.Helper()

	reqBody := dto.CommentImageRequest{
		ImageID: imageID,
		Comment: comment,
	}
	bodyJson, err := json.Marshal(reqBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images/_comment", bytes.NewReader(bodyJson))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)
}

func likeImage(t *testing.T, token string, imageID int64) {
	t.Helper()

	reqBody := dto.LikeImageRequest{
		ImageID: imageID,
	}
	bodyJson, err := json.Marshal(reqBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images/_like", bytes.NewReader(bodyJson))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)
}

func requireNil(t *testing.T, f func() error) {
	require.Nil(t, f())
}

func ClearAll() {
	ClearComments()
	ClearLikes()
	ClearImages()
	ClearUsers()
	redisClient.FlushAll(context.Background())
}

func ClearUsers() {
	err := db.Unscoped().Where("id is not null").Delete(&entity.User{}).Error
	if err != nil {
		x.Logger.Panicf("Failed clear user data : %+v", err)
	}
}

func ClearImages() {
	err := db.Unscoped().Where("id is not null").Delete(&entity.Image{}).Error
	if err != nil {
		x.Logger.Panicf("Failed clear image data : %+v", err)
	}
}

func ClearLikes() {
	err := db.Unscoped().Where("id is not null").Delete(&entity.Like{}).Error
	if err != nil {
		x.Logger.Panicf("Failed clear like data : %+v", err)
	}
}

func ClearComments() {
	err := db.Unscoped().Where("id is not null").Delete(&entity.Comment{}).Error
	if err != nil {
		x.Logger.Panicf("Failed clear comment data : %+v", err)
	}
}

func GetFirstUser(t *testing.T) *entity.User {
	user := &entity.User{}
	err := db.First(user).Error
	require.Nil(t, err)
	return user
}
