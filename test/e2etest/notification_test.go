package e2etest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/http/response"
	"github.com/stretchr/testify/require"
)

func TestFollowCreatesNotification(t *testing.T) {
	ClearAll()

	// User A will be notified, User B performs the follow.
	tokenA := registerAndLoginUser(t, "notif_a", "password", "Notif A")
	tokenB := registerAndLoginUser(t, "notif_b", "password", "Notif B")

	userA := getUserByUsername(t, "notif_a")
	userB := getUserByUsername(t, "notif_b")

	// Before any activity, User A has nothing unread.
	require.Equal(t, int64(0), getUnreadCount(t, tokenA))

	// User B follows User A, which emits an outbox event consumed asynchronously.
	followUser(t, tokenB, userA.ID)

	// The notification is created by a Kafka consumer, so poll until it lands.
	notifications := waitForNotification(t, tokenA, entity.NotificationTypeFollow)
	require.Len(t, notifications.Data, 1)

	notification := notifications.Data[0]
	require.Equal(t, userB.ID, notification.ActorID)
	require.Equal(t, "Notif B", notification.ActorName)
	require.Equal(t, entity.NotificationTypeFollow, notification.Type)
	require.False(t, notification.IsRead)

	// The unread count reflects the single follow notification.
	require.Equal(t, int64(1), getUnreadCount(t, tokenA))
}

func TestNotificationMarkReadAndMarkAllRead(t *testing.T) {
	ClearAll()

	tokenA := registerAndLoginUser(t, "notif_read_a", "password", "Notif Read A")
	tokenB := registerAndLoginUser(t, "notif_read_b", "password", "Notif Read B")

	userA := getUserByUsername(t, "notif_read_a")

	followUser(t, tokenB, userA.ID)
	waitForNotification(t, tokenA, entity.NotificationTypeFollow)

	notifications := listNotifications(t, tokenA, 10)
	require.Len(t, notifications.Data, 1)
	notificationID := notifications.Data[0].ID

	// Marking a single notification read drops the unread count to zero.
	markNotificationRead(t, tokenA, notificationID)
	require.Equal(t, int64(0), getUnreadCount(t, tokenA))

	// It is still listed, now flagged as read.
	notifications = listNotifications(t, tokenA, 10)
	require.Len(t, notifications.Data, 1)
	require.True(t, notifications.Data[0].IsRead)

	// A second notification, then mark-all-read clears the unread count again.
	tokenC := registerAndLoginUser(t, "notif_read_c", "password", "Notif Read C")
	followUser(t, tokenC, userA.ID)
	waitForUnreadCount(t, tokenA, 1)

	markAllNotificationsRead(t, tokenA)
	require.Equal(t, int64(0), getUnreadCount(t, tokenA))
}

// getUserByUsername loads a user row directly so the test can reference the
// numeric ID the API expects (the HTTP layer exposes profiles by ID only).
func getUserByUsername(t *testing.T, username string) *entity.User {
	t.Helper()

	user := &entity.User{}
	err := db.Where("username = ?", username).First(user).Error
	require.Nil(t, err)
	return user
}

// listNotifications fetches the current user's notification page over HTTP.
func listNotifications(t *testing.T, token string, limit int) dto.NotificationsResponse {
	t.Helper()

	url := "http://127.0.0.1:3000/api/notifications?limit=" + strconv.Itoa(limit)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.NotificationsResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	return responseBody.Data
}

// waitForNotification polls the list endpoint until at least minCount
// notifications of the given type appear, since creation is asynchronous.
func waitForNotification(t *testing.T, token, notificationType string) dto.NotificationsResponse {
	t.Helper()

	deadline := time.Now().Add(15 * time.Second)
	for {
		notifications := listNotifications(t, token, 10)
		for _, notification := range notifications.Data {
			if notification.Type == notificationType {
				return notifications
			}
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %q notification to appear", notificationType)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func getUnreadCount(t *testing.T, token string) int64 {
	t.Helper()

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/api/notifications/_unread_count", nil)
	require.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	require.Nil(t, err)

	responseBody := &response.WebResponse[dto.UnreadCountResponse]{}
	err = json.Unmarshal(bytes, responseBody)
	require.Nil(t, err)

	return responseBody.Data.UnreadCount
}

func waitForUnreadCount(t *testing.T, token string, want int64) {
	t.Helper()

	deadline := time.Now().Add(15 * time.Second)
	for {
		if getUnreadCount(t, token) == want {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for unread count to reach %d", want)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func markNotificationRead(t *testing.T, token string, notificationID int64) {
	t.Helper()

	reqBody := dto.MarkNotificationReadRequest{NotificationID: notificationID}
	bodyJson, err := json.Marshal(reqBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/notifications/_read", bytes.NewReader(bodyJson))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)
}

func markAllNotificationsRead(t *testing.T, token string) {
	t.Helper()

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/notifications/_read_all", nil)
	require.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)
}
