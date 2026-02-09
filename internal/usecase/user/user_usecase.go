package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/storage"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseUser.go -pkg=mock . UserUsecase

type UserUsecase interface {
	Verify(ctx context.Context, req *model.VerifyUserRequest) (*model.UserAuth, error)
	Create(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error)
	Login(ctx context.Context, req *model.LoginUserRequest) (*model.UserLoginResponse, error)
	Current(ctx context.Context, req *model.GetUserRequest) (*model.UserResponse, error)
	Update(ctx context.Context, req *model.UpdateUserRequest) (*model.UserResponse, error)
	Follow(ctx context.Context, req *model.FollowUserRequest) error
	NotifyUserBeingFollowed(ctx context.Context, req *model.NotifyUserBeingFollowedRequest) error
	BatchUpdateUserFollowStats(ctx context.Context, req *model.BatchUpdateUserFollowStatsRequest) error
}

var _ UserUsecase = &UserUsecaseImpl{}

type UserUsecaseImpl struct {
	Config *viper.Viper
	DB     *gorm.DB

	// repository
	UserRepository   repository.UserRepository
	FollowRepository repository.FollowRepository

	// producer
	UserProducer  messaging.UserProducer
	NotifProducer messaging.NotifProducer

	// client
	S3Client storage.S3Client
}

func NewUserUsecase(
	cfg *viper.Viper,
	db *gorm.DB,

	// repository
	userRepository repository.UserRepository,
	FollowRepository repository.FollowRepository,

	// producer
	userProducer messaging.UserProducer,
	NotifProducer messaging.NotifProducer,

	// client
	s3Client storage.S3Client,
) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		Config: cfg,
		DB:     db,

		// repository
		UserRepository:   userRepository,
		FollowRepository: FollowRepository,

		// producer
		UserProducer:  userProducer,
		NotifProducer: NotifProducer,

		// client
		S3Client: s3Client,
	}
}
