package user

import (
	"context"
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/gateway/rest"
	"golang-clean-architecture/internal/model"
	"golang-clean-architecture/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserUseCase interface {
	Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.Auth, error)
	Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserResponse, error)
	Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error)
	Logout(ctx context.Context, req *model.LogoutUserRequest) (bool, error)
	Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error)
}

var _ UserUseCase = &UserUseCaseImpl{}

type UserUseCaseImpl struct {
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

func NewUserUseCase(
	db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,

	// repository
	userRepository repository.UserRepository,

	// producer
	userProducer messaging.UserProducer,

	// client
	s3Client rest.S3Client,
	slackClient rest.SlackClient,
) *UserUseCaseImpl {
	return &UserUseCaseImpl{
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
