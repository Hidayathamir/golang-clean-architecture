package converter

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
)

func DtoEmailVerificationEventToDtoSendVerificationEmailRequest(event dto.EmailVerificationEvent, req *dto.SendVerificationEmailRequest) {
	req.To = event.Email
	req.Name = event.Name
	req.Token = event.Token
}
