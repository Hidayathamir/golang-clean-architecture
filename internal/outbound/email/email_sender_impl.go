package email

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
)

var _ Sender = &SenderImpl{}

type SenderImpl struct {
	Cfg *config.Config
}

func NewSender(cfg *config.Config) *SenderImpl {
	return &SenderImpl{
		Cfg: cfg,
	}
}

// SendVerificationEmail delivers a verification link over SMTP. Mailpit listens
// on the configured host:port and captures the message; no auth/TLS is needed
// locally, which keeps this example free of credentials.
func (s *SenderImpl) SendVerificationEmail(ctx context.Context, req dto.SendVerificationEmailRequest) error {
	if err := ctx.Err(); err != nil {
		return errkit.AddFuncName(err, "email.(*SenderImpl).SendVerificationEmail")
	}

	link := fmt.Sprintf("%s?token=%s", s.Cfg.GetEmailVerificationURL(), req.Token)

	from := s.Cfg.GetEmailFrom()
	addr := fmt.Sprintf("%s:%d", s.Cfg.GetEmailSMTPHost(), s.Cfg.GetEmailSMTPPort())

	headers := []string{
		fmt.Sprintf("From: %s", from),
		fmt.Sprintf("To: %s", req.To),
		"Subject: Verify your email address",
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=\"utf-8\"",
	}

	body := fmt.Sprintf(
		"Hi %s,\r\n\r\nConfirm your email address by opening this link:\r\n\r\n%s\r\n\r\nIf you did not sign up, you can ignore this message.\r\n",
		req.Name, link,
	)

	message := strings.Join(headers, "\r\n") + "\r\n\r\n" + body

	err := smtp.SendMail(addr, nil, from, []string{req.To}, []byte(message))
	if err != nil {
		err = errkit.SetCode(err, 502)
		return errkit.AddFuncName(err, "email.(*SenderImpl).SendVerificationEmail")
	}

	return nil
}
