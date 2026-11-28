package e2etest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// mailpitAPIURL is where the dev Mailpit instance exposes its HTTP API. It
// mirrors the SMTP host the app is configured with (see settings/config.json).
const mailpitAPIURL = "http://127.0.0.1:8025/api/v1"

var verificationTokenRe = regexp.MustCompile(`token=([0-9a-f]+)`)

type mailpitMessage struct {
	ID      string `json:"ID"`
	Subject string `json:"Subject"`
	Snippet string `json:"Snippet"`
	To      []struct {
		Address string `json:"Address"`
	} `json:"To"`
}

type mailpitMessagesResponse struct {
	Total    int              `json:"total"`
	Messages []mailpitMessage `json:"messages"`
}

// purgeMailpit empties the Mailpit inbox so a test only sees messages it sent.
func purgeMailpit(t *testing.T) {
	t.Helper()

	req, err := http.NewRequest(http.MethodDelete, mailpitAPIURL+"/messages", nil)
	require.Nil(t, err)

	res, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)
}

// waitForVerificationToken polls Mailpit until a verification email addressed
// to `to` arrives, then returns the token embedded in its link. It polls rather
// than sleeps because delivery is asynchronous: register -> outbox -> Kafka ->
// consumer -> SMTP.
func waitForVerificationToken(t *testing.T, to string) string {
	t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	for {
		token, found := findVerificationToken(to)
		if found {
			return token
		}
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for verification email to %s", to)
		}
		time.Sleep(250 * time.Millisecond)
	}
}

func findVerificationToken(to string) (string, bool) {
	res, err := http.Get(mailpitAPIURL + "/messages")
	if err != nil {
		return "", false
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", false
	}

	response := mailpitMessagesResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", false
	}

	for _, m := range response.Messages {
		if m.Subject != "Verify your email address" {
			continue
		}
		matchedRecipient := false
		for _, r := range m.To {
			if r.Address == to {
				matchedRecipient = true
				break
			}
		}
		if !matchedRecipient {
			continue
		}
		match := verificationTokenRe.FindStringSubmatch(m.Snippet)
		if len(match) == 2 {
			return match[1], true
		}
	}

	return "", false
}

// verifyEmail calls the verification endpoint with the given token and returns
// the HTTP status code.
func verifyEmail(t *testing.T, token string) int {
	t.Helper()

	url := fmt.Sprintf("http://127.0.0.1:3000/api/users/_verify_email?token=%s", token)
	res, err := http.Get(url)
	require.Nil(t, err)
	defer requireNil(t, res.Body.Close)

	return res.StatusCode
}
