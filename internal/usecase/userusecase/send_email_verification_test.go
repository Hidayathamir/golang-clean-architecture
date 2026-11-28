package userusecase_test

import (
	"context"
	"testing"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/mock"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/userusecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserUsecaseImpl_SendEmailVerification_Success(t *testing.T) {
	emailSender := &mock.SenderMock{}
	u := &userusecase.UserUsecaseImpl{
		EmailSender: emailSender,
	}

	// ------------------------------------------------------- //

	var sent dto.SendVerificationEmailRequest
	emailSender.SendVerificationEmailFunc = func(ctx context.Context, req dto.SendVerificationEmailRequest) error {
		sent = req
		return nil
	}

	// ------------------------------------------------------- //

	err := u.SendEmailVerification(context.Background(), dto.EmailVerificationEvent{
		UserID: 7,
		Email:  "user1@example.com",
		Name:   "name1",
		Token:  "token-123",
	})

	// ------------------------------------------------------- //

	require.Nil(t, err)
	require.Equal(t, "user1@example.com", sent.To)
	require.Equal(t, "name1", sent.Name)
	require.Equal(t, "token-123", sent.Token)
}

func TestUserUsecaseImpl_SendEmailVerification_Fail_ValidateStruct(t *testing.T) {
	u := &userusecase.UserUsecaseImpl{
		EmailSender: &mock.SenderMock{},
	}

	err := u.SendEmailVerification(context.Background(), dto.EmailVerificationEvent{})

	require.NotNil(t, err)
}

func TestUserUsecaseImpl_SendEmailVerification_Fail_Sender(t *testing.T) {
	emailSender := &mock.SenderMock{}
	u := &userusecase.UserUsecaseImpl{
		EmailSender: emailSender,
	}

	emailSender.SendVerificationEmailFunc = func(ctx context.Context, req dto.SendVerificationEmailRequest) error {
		return assert.AnError
	}

	err := u.SendEmailVerification(context.Background(), dto.EmailVerificationEvent{
		UserID: 7,
		Email:  "user1@example.com",
		Token:  "token-123",
	})

	require.NotNil(t, err)
}
