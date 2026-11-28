package e2etest

import (
	"net/http"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/column"
	"github.com/stretchr/testify/require"
)

// TestEmailVerificationFlow covers the whole async path: a registration writes
// an outbox event, the worker turns it into an email (captured by Mailpit), and
// the link's token confirms the account.
func TestEmailVerificationFlow(t *testing.T) {
	ClearAll()
	purgeMailpit(t)

	// Register a user; the verification email is sent asynchronously.
	registerUser(t, "verifyuser", "secret123", "Verify User")

	user := findUserByUsername(t, "verifyuser")
	require.False(t, user.IsEmailVerified(), "user should start unverified")

	// Wait for the email to land and pull the token out of its link.
	token := waitForVerificationToken(t, "verifyuser@example.com")
	require.NotEmpty(t, token)

	// Confirm the account.
	require.Equal(t, http.StatusOK, verifyEmail(t, token))

	verified := findUserByUsername(t, "verifyuser")
	require.True(t, verified.IsEmailVerified(), "email_verified_at should be set")
}

// TestEmailVerificationInvalidToken verifies an unknown token is rejected.
func TestEmailVerificationInvalidToken(t *testing.T) {
	ClearAll()

	require.Equal(t, http.StatusNotFound, verifyEmail(t, "0000000000000000"))
}

// TestEmailVerificationIsIdempotent verifies confirming twice is a no-op success.
func TestEmailVerificationIsIdempotent(t *testing.T) {
	ClearAll()
	purgeMailpit(t)

	registerUser(t, "verifytwice", "secret123", "Verify Twice")

	token := waitForVerificationToken(t, "verifytwice@example.com")

	require.Equal(t, http.StatusOK, verifyEmail(t, token))
	require.Equal(t, http.StatusOK, verifyEmail(t, token))

	user := findUserByUsername(t, "verifytwice")
	require.True(t, user.IsEmailVerified())
}

func findUserByUsername(t *testing.T, username string) *entity.User {
	t.Helper()

	user := &entity.User{}
	where, arg := column.Username.Eq(username)
	err := db.Where(where, arg).Take(user).Error
	require.Nil(t, err)

	return user
}
