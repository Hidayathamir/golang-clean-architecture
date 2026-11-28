package email

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
)

//go:generate moq -out=../../mock/MockEmailSender.go -pkg=mock . Sender

// Sender is the outbound email boundary. Swapping providers (SMTP, SES, a
// transactional API) means a new implementation of this interface only.
type Sender interface {
	SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error
}
