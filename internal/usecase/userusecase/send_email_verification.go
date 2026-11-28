package userusecase

import (
	"context"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

// SendEmailVerification is the Kafka consumer entry point for the
// email.verification topic. It maps the event to a sender request; the actual
// provider call lives behind the outbound email.Sender interface.
func (u *UserUsecaseImpl) SendEmailVerification(ctx context.Context, req dto.EmailVerificationEvent) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).SendEmailVerification")
	}

	sendReq := dto.SendVerificationEmailRequest{}
	converter.DtoEmailVerificationEventToDtoSendVerificationEmailRequest(req, &sendReq)

	err = u.EmailSender.SendVerificationEmail(ctx, sendReq)
	if err != nil {
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).SendEmailVerification")
	}

	return nil
}
