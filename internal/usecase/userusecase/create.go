package userusecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/Hidayathamir/golang-clean-architecture/internal/converter"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/validatorkit"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (u *UserUsecaseImpl) Create(ctx context.Context, req dto.RegisterUserRequest) (dto.UserResponse, error) {
	err := validatorkit.Validate.Struct(&req)
	if err != nil {
		err = errkit.SetCode(err, http.StatusBadRequest)
		return dto.UserResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Create")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Create")
	}

	token, err := generateVerificationToken()
	if err != nil {
		return dto.UserResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Create")
	}

	user := entity.User{}
	converter.DtoRegisterUserRequestToEntityUser(req, &user, string(password))
	user.VerificationToken = token

	// The user row and its verification event commit together, so registration
	// never succeeds without an email on the way, and vice versa.
	err = u.DB.Transaction(func(tx *gorm.DB) error {
		err := u.UserRepository.Create(ctx, tx, &user)
		if err != nil {
			return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Create")
		}

		event := dto.EmailVerificationEvent{
			UserID: user.ID,
			Email:  user.Email,
			Name:   user.Name,
			Token:  user.VerificationToken,
		}
		err = u.EmailProducer.SendVerificationEmail(ctx, tx, &event)
		if err != nil {
			return errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Create")
		}
		return nil
	})
	if err != nil {
		return dto.UserResponse{}, errkit.AddFuncName(err, "userusecase.(*UserUsecaseImpl).Create")
	}

	res := dto.UserResponse{}
	converter.EntityUserToDtoUserResponse(user, &res)

	return res, nil
}

// generateVerificationToken returns a 32-byte random token, hex encoded. The
// value is unguessable and short enough to travel in a URL.
func generateVerificationToken() (string, error) {
	buf := make([]byte, 32)
	_, err := rand.Read(buf)
	if err != nil {
		return "", errkit.AddFuncName(err, "userusecase.generateVerificationToken")
	}
	return hex.EncodeToString(buf), nil
}
