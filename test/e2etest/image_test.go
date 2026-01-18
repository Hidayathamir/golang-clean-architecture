package e2etest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/response"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func uploadImage(t *testing.T, token string) int64 {
	t.Helper()
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "test.jpg")
	assert.Nil(t, err)
	_, err = part.Write([]byte("dummy image content"))
	assert.Nil(t, err)
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images", body)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	bytesBody, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.ImageResponse])
	err = json.Unmarshal(bytesBody, responseBody)
	assert.Nil(t, err)

	return responseBody.Data.ID
}

func TestUploadImage(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "test.jpg")
	assert.Nil(t, err)
	_, err = part.Write([]byte("dummy image content"))
	assert.Nil(t, err)
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images", body)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	bytesBody, err := io.ReadAll(res.Body)
	assert.Nil(t, err)

	responseBody := new(response.WebResponse[model.ImageResponse])
	err = json.Unmarshal(bytesBody, responseBody)
	assert.Nil(t, err)

	assert.NotZero(t, responseBody.Data.ID)
	assert.NotEmpty(t, responseBody.Data.URL)

	var count int64
	err = db.Model(&entity.Image{}).Where("id = ?", responseBody.Data.ID).Count(&count).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func TestLikeImage(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)
	imageID := uploadImage(t, token)

	reqBody := model.LikeImageRequest{
		ImageID: imageID,
	}
	bodyJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images/_like", bytes.NewReader(bodyJson))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var count int64
	err = db.Model(&entity.Like{}).Where("image_id = ?", imageID).Count(&count).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func TestCommentImage(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)
	imageID := uploadImage(t, token)

	reqBody := model.CommentImageRequest{
		ImageID: imageID,
		Comment: "Nice!",
	}
	bodyJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images/_comment", bytes.NewReader(bodyJson))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	var count int64
	err = db.Model(&entity.Comment{}).Where("image_id = ? AND comment = ?", imageID, "Nice!").Count(&count).Error
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func TestGetLikes(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)
	imageID := uploadImage(t, token)

	// Like it first
	reqBody := model.LikeImageRequest{
		ImageID: imageID,
	}
	bodyJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)
	reqLike, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images/_like", bytes.NewReader(bodyJson))
	assert.Nil(t, err)
	reqLike.Header.Set("Content-Type", "application/json")
	reqLike.Header.Set("Authorization", bearerToken(token))
	resLike, err := http.DefaultClient.Do(reqLike)
	assert.Nil(t, err)
	resLike.Body.Close()

	// Get Likes
	url := fmt.Sprintf("http://127.0.0.1:3000/api/images/%d/likes", imageID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	respBody := new(response.WebResponse[model.LikeResponseList])
	err = json.NewDecoder(res.Body).Decode(respBody)
	assert.Nil(t, err)
	assert.NotEmpty(t, respBody.Data)
	assert.Equal(t, imageID, respBody.Data[0].ImageID)
}

func TestGetComments(t *testing.T) {
	ClearAll()
	token := registerAndLoginDefaultUser(t)
	imageID := uploadImage(t, token)

	// Comment first
	reqBody := model.CommentImageRequest{
		ImageID: imageID,
		Comment: "Wow",
	}
	bodyJson, err := json.Marshal(reqBody)
	assert.Nil(t, err)
	reqComm, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/images/_comment", bytes.NewReader(bodyJson))
	assert.Nil(t, err)
	reqComm.Header.Set("Content-Type", "application/json")
	reqComm.Header.Set("Authorization", bearerToken(token))
	resComm, err := http.DefaultClient.Do(reqComm)
	assert.Nil(t, err)
	resComm.Body.Close()

	// Get Comments
	url := fmt.Sprintf("http://127.0.0.1:3000/api/images/%d/comments", imageID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	respBody := new(response.WebResponse[model.CommentResponseList])
	err = json.NewDecoder(res.Body).Decode(respBody)
	assert.Nil(t, err)
	assert.NotEmpty(t, respBody.Data)
	assert.Equal(t, "Wow", respBody.Data[0].Comment)
}
