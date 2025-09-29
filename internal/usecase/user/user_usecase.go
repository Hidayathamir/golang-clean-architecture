package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/UserUsecase.go -pkg=mock . UserUsecase

type UserUsecase interface {
	Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error)
	Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error)
	Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error)
	Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error)
	Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error)
}

var _ UserUsecase = &UserUsecaseImpl{}

type UserUsecaseImpl struct {
	DB       *gorm.DB
	Log      *logrus.Logger
	Validate *validator.Validate

	// repository
	UserRepository repository.UserRepository

	// producer
	UserProducer messaging.UserProducer

	// client
	S3Client    rest.S3Client
	SlackClient rest.SlackClient
}

func NewUserUsecase(
	db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,

	// repository
	userRepository repository.UserRepository,

	// producer
	userProducer messaging.UserProducer,

	// client
	s3Client rest.S3Client,
	slackClient rest.SlackClient,
) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		DB:       db,
		Log:      logger,
		Validate: validate,

		// repository
		UserRepository: userRepository,

		// producer
		UserProducer: userProducer,

		// client
		S3Client:    s3Client,
		SlackClient: slackClient,
	}
}
