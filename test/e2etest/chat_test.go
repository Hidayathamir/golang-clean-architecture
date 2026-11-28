package e2etest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/http/response"
	"github.com/stretchr/testify/require"
)

func TestChatLifecycle(t *testing.T) {
	ClearAll()

	tokenA := registerAndLoginUser(t, "chat_a", "password", "Chat A")
	tokenB := registerAndLoginUser(t, "chat_b", "password", "Chat B")

	userA := getUserByUsername(t, "chat_a")
	userB := getUserByUsername(t, "chat_b")

	// B starts with no unread chats.
	require.Equal(t, int64(0), getUnreadChatCount(t, tokenB))

	// A sends B a chat; it is stored synchronously.
	sendChat(t, tokenA, userB.ID, "hello there")

	require.Equal(t, int64(1), getUnreadChatCount(t, tokenB))

	// B lists conversations and sees A as the peer with A's message.
	conversations := listConversations(t, tokenB)
	require.Len(t, conversations.Data, 1)
	require.Equal(t, userA.ID, conversations.Data[0].PeerID)
	require.Equal(t, "hello there", conversations.Data[0].LastChat.Content)

	// B reads the thread, then marks it read.
	thread := getThread(t, tokenB, userA.ID)
	require.Len(t, thread.Data, 1)
	markChatRead(t, tokenB, userA.ID)
	require.Equal(t, int64(0), getUnreadChatCount(t, tokenB))
}

func sendChat(t *testing.T, token string, recipientID int64, content string) {
	t.Helper()
	reqBody := dto.SendChatRequest{RecipientID: recipientID, Content: content}
	body, err := json.Marshal(reqBody)
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/chats", bytes.NewReader(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)
}

func listConversations(t *testing.T, token string) dto.ConversationsResponse {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/api/chats/conversations", nil)
	require.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)
	raw, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	body := &response.WebResponse[dto.ConversationsResponse]{}
	err = json.Unmarshal(raw, body)
	require.Nil(t, err)
	return body.Data
}

func getThread(t *testing.T, token string, peerID int64) dto.ThreadResponse {
	t.Helper()
	url := "http://127.0.0.1:3000/api/chats/" + strconv.FormatInt(peerID, 10)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)
	raw, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	body := &response.WebResponse[dto.ThreadResponse]{}
	err = json.Unmarshal(raw, body)
	require.Nil(t, err)
	return body.Data
}

func markChatRead(t *testing.T, token string, peerID int64) {
	t.Helper()
	body, err := json.Marshal(dto.MarkChatReadRequest{PeerID: peerID})
	require.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3000/api/chats/_read", bytes.NewReader(body))
	require.Nil(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)
}

func getUnreadChatCount(t *testing.T, token string) int64 {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/api/chats/_unread_count", nil)
	require.Nil(t, err)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", bearerToken(token))

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	require.Equal(t, http.StatusOK, res.StatusCode)
	raw, err := io.ReadAll(res.Body)
	require.Nil(t, err)
	body := &response.WebResponse[dto.UnreadChatCountResponse]{}
	err = json.Unmarshal(raw, body)
	require.Nil(t, err)
	return body.Data.UnreadCount
}
