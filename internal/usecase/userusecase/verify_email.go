package userusecase

import (
	"context"
	"net/http"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/logkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
)

func (u *UserUsecaseImpl) VerifyEmail(ctx context.Context, req dto.VerifyEmailRequest) error {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).VerifyEmail")
	}

	user := entity.User{}
	err = u.UserRepository.FindByVerificationToken(ctx, u.DB, &user, req.Token)
	if err != nil {
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).VerifyEmail")
	}

	// Idempotent: a second confirm with the same link is a no-op success.
	if user.IsEmailVerified() {
		return nil
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	err = u.UserRepository.Update(ctx, u.DB, &user)
	if err != nil {
		return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).VerifyEmail")
	}

	// The cached auth view is now stale (it does not carry email state anyway);
	// drop it so the next request reloads the row.
	err = u.UserCache.Delete(ctx, user.ID)
	logkit.LogIfErr(err)

	return nil
}
