package e2etest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/http/response"
	"github.com/stretchr/testify/require"
)

// presignImage asks the server for a presigned S3 PUT URL for the given filename.
func presignImage(t *testing.T, token, filename string) dto.PresignUploadImageResponse {
	t.Helper()

	reqBody := dto.PresignUploadImageRequest{
		Filename:    filename,
		ContentType: "image/jpeg",
	}
	bodyJson, err := json.Marshal(reqBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/images", bytes.NewReader(bodyJson))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)
	require.Equal(t, http.StatusOK, res.StatusCode)

	respBody := &response.WebResponse[dto.PresignUploadImageResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(respBody))
	require.NotEmpty(t, respBody.Data.UploadURL)
	require.NotEmpty(t, respBody.Data.Key)
	return respBody.Data
}

// putToPresignedURL uploads the image bytes directly to S3 via the presigned URL.
func putToPresignedURL(t *testing.T, uploadURL, contentType string, content []byte) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPut, uploadURL, bytes.NewReader(content))
	require.Nil(t, err)
	req.Header.Set("Content-Type", contentType)

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)
	require.Equal(t, http.StatusOK, res.StatusCode)
}

// confirmImage finalizes the upload and returns the new image id.
func confirmImage(t *testing.T, token, key, caption string) int64 {
	t.Helper()

	reqBody := dto.ConfirmUploadImageRequest{
		Key:     key,
		Caption: caption,
	}
	bodyJson, err := json.Marshal(reqBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/images/_confirm", bytes.NewReader(bodyJson))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)
	require.Equal(t, http.StatusOK, res.StatusCode)

	respBody := &response.WebResponse[dto.ImageResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(respBody))
	require.NotZero(t, respBody.Data.ID)
	return respBody.Data.ID
}

func uploadImage(t *testing.T, token string) int64 {
	t.Helper()

	presigned := presignImage(t, token, "test.jpg")
	putToPresignedURL(t, presigned.UploadURL, "image/jpeg", []byte("dummy image content"))
	return confirmImage(t, token, presigned.Key, "")
}

func TestUploadImage(t *testing.T) {
	ClearAll()

	// Register and login user
	token := registerAndLoginDefaultUser(t)

	// Presign, upload bytes directly to S3, then confirm
	presigned := presignImage(t, token, "test.jpg")
	require.NotEmpty(t, presigned.URL)
	putToPresignedURL(t, presigned.UploadURL, "image/jpeg", []byte("dummy image content"))
	imageID := confirmImage(t, token, presigned.Key, "My Caption")

	// Verify database record
	image := &entity.Image{}
	err := db.First(image, imageID).Error
	require.Nil(t, err)
	require.Equal(t, "My Caption", image.Caption)
	require.Equal(t, presigned.URL, image.URL)
}

func TestLikeImage(t *testing.T) {
	ClearAll()

	// Register and login user
	token := registerAndLoginDefaultUser(t)

	// Upload image first
	imageID := uploadImage(t, token)

	// Prepare like request
	reqBody := dto.LikeImageRequest{
		ImageID: imageID,
	}
	bodyJson, err := json.Marshal(reqBody)
	require.Nil(t, err)

	// Send like request
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images/_like", bytes.NewReader(bodyJson))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify response status code
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Verify database record
	var count int64
	err = db.Model(&entity.Like{}).Where("image_id = ?", imageID).Count(&count).Error
	require.Nil(t, err)
	require.Equal(t, int64(1), count)
}

func TestCommentImage(t *testing.T) {
	ClearAll()

	// Register and login user
	token := registerAndLoginDefaultUser(t)

	// Upload image first
	imageID := uploadImage(t, token)

	// Prepare comment request
	reqBody := dto.CommentImageRequest{
		ImageID: imageID,
		Comment: "Nice!",
	}
	bodyJson, err := json.Marshal(reqBody)
	require.Nil(t, err)

	// Send comment request
	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images/_comment", bytes.NewReader(bodyJson))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify response status code
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Verify database record
	var count int64
	err = db.Model(&entity.Comment{}).Where("image_id = ? AND comment = ?", imageID, "Nice!").Count(&count).Error
	require.Nil(t, err)
	require.Equal(t, int64(1), count)
}

func TestGetLikes(t *testing.T) {
	ClearAll()

	// Register and login user
	token := registerAndLoginDefaultUser(t)

	// Upload image first
	imageID := uploadImage(t, token)

	// Like it first
	reqBody := dto.LikeImageRequest{
		ImageID: imageID,
	}
	bodyJson, err := json.Marshal(reqBody)
	require.Nil(t, err)
	reqLike, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images/_like", bytes.NewReader(bodyJson))
	require.Nil(t, err)
	reqLike.Header.Set("Content-Type", "application/json")
	reqLike.Header.Set("Authorization", bearerToken(token))
	resLike, err := http.DefaultClient.Do(reqLike)
	require.Nil(t, err)
	require.Nil(t, resLike.Body.Close())

	// Send get likes request
	url := fmt.Sprintf("http://127.0.0.1:3000/api/images/%d/likes", imageID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify response status code
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Verify response body
	respBody := &response.WebResponse[dto.LikesResponse]{}
	err = json.NewDecoder(res.Body).Decode(respBody)
	require.Nil(t, err)
	require.NotEmpty(t, respBody.Data.Data)
	require.Equal(t, imageID, respBody.Data.Data[0].ImageID)
}

func TestGetComments(t *testing.T) {
	ClearAll()

	// Register and login user
	token := registerAndLoginDefaultUser(t)

	// Upload image first
	imageID := uploadImage(t, token)

	// Comment first
	reqBody := dto.CommentImageRequest{
		ImageID: imageID,
		Comment: "Wow",
	}
	bodyJson, err := json.Marshal(reqBody)
	require.Nil(t, err)
	reqComm, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images/_comment", bytes.NewReader(bodyJson))
	require.Nil(t, err)
	reqComm.Header.Set("Content-Type", "application/json")
	reqComm.Header.Set("Authorization", bearerToken(token))
	resComm, err := http.DefaultClient.Do(reqComm)
	require.Nil(t, err)
	require.Nil(t, resComm.Body.Close())

	// Send get comments request
	url := fmt.Sprintf("http://127.0.0.1:3000/api/images/%d/comments", imageID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	// Verify response status code
	require.Equal(t, http.StatusOK, res.StatusCode)

	// Verify response body
	respBody := &response.WebResponse[dto.CommentsResponse]{}
	err = json.NewDecoder(res.Body).Decode(respBody)
	require.Nil(t, err)
	require.NotEmpty(t, respBody.Data.Data)
	require.Equal(t, "Wow", respBody.Data.Data[0].Comment)
}

func TestImageFlow(t *testing.T) {
	ClearAll()

	// Register multiple users
	token1 := registerAndLoginUser(t, "user1", "password", "User One")
	token2 := registerAndLoginUser(t, "user2", "password", "User Two")
	token3 := registerAndLoginUser(t, "user3", "password", "User Three")

	// Get User 1 ID
	user1 := &entity.User{}
	err := db.Where("username = ?", "user1").First(user1).Error
	require.Nil(t, err)

	// User 2 follows User 1
	reqBodyFollow := dto.FollowUserRequest{
		FollowingID: user1.ID,
	}
	bodyJsonFollow, err := json.Marshal(reqBodyFollow)
	require.Nil(t, err)

	reqFollow2, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_follow", bytes.NewReader(bodyJsonFollow))
	require.Nil(t, err)
	reqFollow2.Header.Set("Content-Type", "application/json")
	reqFollow2.Header.Set("Authorization", bearerToken(token2))

	resFollow2, err := http.DefaultClient.Do(reqFollow2)
	require.Nil(t, err)
	require.Nil(t, resFollow2.Body.Close())
	require.Equal(t, http.StatusOK, resFollow2.StatusCode)

	// User 3 follows User 1
	reqFollow3, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/users/_follow", bytes.NewReader(bodyJsonFollow))
	require.Nil(t, err)
	reqFollow3.Header.Set("Content-Type", "application/json")
	reqFollow3.Header.Set("Authorization", bearerToken(token3))

	resFollow3, err := http.DefaultClient.Do(reqFollow3)
	require.Nil(t, err)
	require.Nil(t, resFollow3.Body.Close())
	require.Equal(t, http.StatusOK, resFollow3.StatusCode)

	// User 1 uploads image
	uploadImage(t, token1)
}

func TestMultipleUsersCommentImage(t *testing.T) {
	ClearAll()

	// Register users
	tokenA := registerAndLoginUser(t, "userA", "password", "User A")
	tokenB := registerAndLoginUser(t, "userB", "password", "User B")
	tokenC := registerAndLoginUser(t, "userC", "password", "User C")

	// User A uploads image
	imageID := uploadImage(t, tokenA)

	// User B comments
	commentImage(t, tokenB, imageID, "Comment from B")

	// User C comments
	commentImage(t, tokenC, imageID, "Comment from C")

	// Verify comments
	url := fmt.Sprintf("http://127.0.0.1:3000/api/images/%d/comments", imageID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(tokenA))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)

	respBody := &response.WebResponse[dto.CommentsResponse]{}
	err = json.NewDecoder(res.Body).Decode(respBody)
	require.Nil(t, err)

	require.Len(t, respBody.Data.Data, 2)
}

func TestMultipleUsersLikeImage(t *testing.T) {
	ClearAll()

	// Register users
	tokenA := registerAndLoginUser(t, "userA", "password", "User A")
	tokenB := registerAndLoginUser(t, "userB", "password", "User B")
	tokenC := registerAndLoginUser(t, "userC", "password", "User C")

	// User A uploads image
	imageID := uploadImage(t, tokenA)

	// User B likes
	likeImage(t, tokenB, imageID)

	// User C likes
	likeImage(t, tokenC, imageID)

	// Verify likes
	url := fmt.Sprintf("http://127.0.0.1:3000/api/images/%d/likes", imageID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(tokenA))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)

	respBody := &response.WebResponse[dto.LikesResponse]{}
	err = json.NewDecoder(res.Body).Decode(respBody)
	require.Nil(t, err)

	require.Len(t, respBody.Data.Data, 2)
}
