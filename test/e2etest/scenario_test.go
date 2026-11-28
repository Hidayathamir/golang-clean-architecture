package e2etest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/http/response"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://127.0.0.1:3000"

// getUserID looks up a user's id by username from the database.
func getUserID(t *testing.T, username string) int64 {
	t.Helper()
	user := &entity.User{}
	err := db.Where("username = ?", username).First(user).Error
	require.Nil(t, err)
	return user.ID
}

// uploadImageWithCaption uploads an image with a caption and returns the new image id.
func uploadImageWithCaption(t *testing.T, token, caption string) int64 {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "photo.jpg")
	require.Nil(t, err)
	_, err = part.Write([]byte("dummy image content"))
	require.Nil(t, err)
	require.Nil(t, writer.WriteField("caption", caption))
	require.Nil(t, writer.Close())

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/images", body)
	require.Nil(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
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

// getFeed fetches the current user's feed page.
func getFeed(t *testing.T, token, query string) dto.FeedResponse {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/feed"+query, nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)
	require.Equal(t, http.StatusOK, res.StatusCode)

	respBody := &response.WebResponse[dto.FeedResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(respBody))
	return respBody.Data
}

// getProfile fetches a user's public profile.
func getProfile(t *testing.T, token string, userID int64) dto.ProfileResponse {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/users/%d", baseURL, userID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)
	require.Equal(t, http.StatusOK, res.StatusCode)

	resp := &response.WebResponse[dto.ProfileResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(resp))
	return resp.Data
}

// eventually retries cond for a short while, since follow stats are updated
// asynchronously by a Kafka consumer (a real client would just refresh).
func eventually(t *testing.T, cond func() bool) {
	t.Helper()
	for i := 0; i < 50; i++ {
		if cond() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
	require.Fail(t, "condition not met within timeout")
}

// TestRealUserJourney walks the full Instagram-like flow over HTTP, exactly the
// way a real user (or mobile app) would: register, log in, follow, post, react,
// browse profiles, then unwind (unlike, unfollow, delete).
func TestRealUserJourney(t *testing.T) {
	ClearAll()

	// --- Sign up: Alice and Bob create accounts and log in ---
	aliceToken := registerAndLoginUser(t, "alice", "password", "Alice")
	bobToken := registerAndLoginUser(t, "bob", "password", "Bob")
	aliceID := getUserID(t, "alice")
	bobID := getUserID(t, "bob")

	// --- Social graph: each follows the other ---
	followUser(t, aliceToken, bobID)
	followUser(t, bobToken, aliceID)
	checkFollow(t, aliceID, bobID)
	checkFollow(t, bobID, aliceID)

	// --- Bob posts two photos ---
	firstImageID := uploadImageWithCaption(t, bobToken, "Bob's first post")
	secondImageID := uploadImageWithCaption(t, bobToken, "Bob's second post")

	// --- Alice's feed shows Bob's posts, newest first ---
	feed := getFeed(t, aliceToken, "")
	require.Len(t, feed.Data, 2)
	require.Equal(t, secondImageID, feed.Data[0].ID) // newest first
	require.Equal(t, firstImageID, feed.Data[1].ID)

	// --- Alice opens a single post ---
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/images/%d", baseURL, firstImageID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(aliceToken))
	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	respImage := &response.WebResponse[dto.ImageResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(respImage))
	require.Nil(t, res.Body.Close())
	require.Equal(t, "Bob's first post", respImage.Data.Caption)
	require.Equal(t, bobID, respImage.Data.UserID)

	// --- Alice interacts: likes and comments on the post ---
	likeImage(t, aliceToken, firstImageID)
	commentImage(t, aliceToken, firstImageID, "Love this!")

	// Likes list reflects the interaction.
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/images/%d/likes", baseURL, firstImageID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(aliceToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	respLikes := &response.WebResponse[dto.LikesResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(respLikes))
	require.Nil(t, res.Body.Close())
	require.Len(t, respLikes.Data.Data, 1)
	require.Equal(t, aliceID, respLikes.Data.Data[0].UserID)

	// Comments list reflects the comment.
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/images/%d/comments", baseURL, firstImageID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(aliceToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	respComments := &response.WebResponse[dto.CommentsResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(respComments))
	require.Nil(t, res.Body.Close())
	require.Len(t, respComments.Data.Data, 1)
	require.Equal(t, "Love this!", respComments.Data.Data[0].Comment)

	// --- Client fetches Bob's profile and his post grid ---
	// Follower counts are updated asynchronously, so poll until they settle.
	eventually(t, func() bool {
		return getProfile(t, aliceToken, bobID).FollowerCount == 1
	})
	profile := getProfile(t, aliceToken, bobID)
	require.Equal(t, "bob", profile.Username)
	require.Equal(t, 1, profile.FollowerCount) // only alice follows bob

	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/users/%d/images", baseURL, bobID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(aliceToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	respGrid := &response.WebResponse[dto.UserImagesResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(respGrid))
	require.Nil(t, res.Body.Close())
	require.Len(t, respGrid.Data.Data, 2)

	// --- Client views Bob's followers and followings lists ---
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/users/%d/followers", baseURL, bobID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(aliceToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	respFollowers := &response.WebResponse[dto.FollowersResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(respFollowers))
	require.Nil(t, res.Body.Close())
	require.Len(t, respFollowers.Data.Data, 1)
	require.Equal(t, "alice", respFollowers.Data.Data[0].Username)

	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/users/%d/following", baseURL, aliceID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(aliceToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	respFollowing := &response.WebResponse[dto.FollowingsResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(respFollowing))
	require.Nil(t, res.Body.Close())
	require.Len(t, respFollowing.Data.Data, 1)
	require.Equal(t, "bob", respFollowing.Data.Data[0].Username)

	// --- Bob edits his first post caption ---
	updateBody, err := json.Marshal(dto.UpdateImageRequest{Caption: "Bob's edited caption"})
	require.Nil(t, err)
	req, err = http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/images/%d", baseURL, firstImageID), bytes.NewReader(updateBody))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(bobToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Nil(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)

	// The edit is visible on the post.
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/images/%d", baseURL, firstImageID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(aliceToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	respImage = &response.WebResponse[dto.ImageResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(respImage))
	require.Nil(t, res.Body.Close())
	require.Equal(t, "Bob's edited caption", respImage.Data.Caption)

	// --- Unwind: Alice unlikes, then unfollows Bob ---
	req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/images/%d/like", baseURL, firstImageID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(aliceToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Nil(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)

	req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/users/%d/follow", baseURL, bobID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(aliceToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Nil(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)

	// The follow relationship is gone and the like row is soft-deleted.
	var followCount int64
	require.Nil(t, db.Model(&entity.Follow{}).Where("follower_id = ? AND following_id = ?", aliceID, bobID).Count(&followCount).Error)
	require.Equal(t, int64(0), followCount)

	// Alice's feed is now empty because she no longer follows Bob.
	feed = getFeed(t, aliceToken, "")
	require.Empty(t, feed.Data)

	// --- Bob deletes his remaining post; his grid shrinks ---
	req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/images/%d", baseURL, secondImageID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(bobToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Nil(t, res.Body.Close())
	require.Equal(t, http.StatusOK, res.StatusCode)

	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/users/%d/images", baseURL, bobID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(bobToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	respGrid = &response.WebResponse[dto.UserImagesResponse]{}
	require.Nil(t, json.NewDecoder(res.Body).Decode(respGrid))
	require.Nil(t, res.Body.Close())
	require.Len(t, respGrid.Data.Data, 1)
	require.Equal(t, firstImageID, respGrid.Data.Data[0].ID)
}

// TestRealUserOwnershipRules verifies that a user cannot act on someone else's
// content, the way the API enforces it for a real client.
func TestRealUserOwnershipRules(t *testing.T) {
	ClearAll()

	ownerToken := registerAndLoginUser(t, "owner", "password", "Owner")
	strangerToken := registerAndLoginUser(t, "stranger", "password", "Stranger")
	ownerID := getUserID(t, "owner")

	imageID := uploadImageWithCaption(t, ownerToken, "Owner's post")

	// Stranger tries to edit the owner's post -> 403.
	updateBody, err := json.Marshal(dto.UpdateImageRequest{Caption: "hacked"})
	require.Nil(t, err)
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/api/images/%d", baseURL, imageID), bytes.NewReader(updateBody))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken(strangerToken))
	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	_, _ = io.Copy(io.Discard, res.Body)
	require.Nil(t, res.Body.Close())
	require.Equal(t, http.StatusForbidden, res.StatusCode)

	// Stranger tries to delete the owner's post -> 403.
	req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/images/%d", baseURL, imageID), nil)
	require.Nil(t, err)
	req.Header.Set("Authorization", bearerToken(strangerToken))
	res, err = http.DefaultClient.Do(req)
	require.Nil(t, err)
	_, _ = io.Copy(io.Discard, res.Body)
	require.Nil(t, res.Body.Close())
	require.Equal(t, http.StatusForbidden, res.StatusCode)

	// The post is untouched and still owned by the owner.
	image := &entity.Image{}
	require.Nil(t, db.First(image, imageID).Error)
	require.Equal(t, ownerID, image.UserID)
	require.Equal(t, "Owner's post", image.Caption)
}
