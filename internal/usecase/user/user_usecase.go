package user

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/storage"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseUser.go -pkg=mock . UserUsecase

type UserUsecase interface {
	Verify(ctx context.Context, req *dto.VerifyUserRequest) (*dto.UserAuth, error)
	Create(ctx context.Context, req *dto.RegisterUserRequest) (*dto.UserResponse, error)
	Login(ctx context.Context, req *dto.LoginUserRequest) (*dto.UserLoginResponse, error)
	Current(ctx context.Context, req *dto.GetUserRequest) (*dto.UserResponse, error)
	Update(ctx context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	Follow(ctx context.Context, req *dto.FollowUserRequest) error
	NotifyUserBeingFollowed(ctx context.Context, req *dto.NotifyUserBeingFollowedRequest) error
	BatchUpdateUserFollowStats(ctx context.Context, req *dto.BatchUpdateUserFollowStatsRequest) error
}

var _ UserUsecase = &UserUsecaseImpl{}

type UserUsecaseImpl struct {
	Config *config.Config
	DB     *gorm.DB

	// repository
	UserRepository     repository.UserRepository
	UserStatRepository repository.UserStatRepository
	FollowRepository   repository.FollowRepository

	// producer
	UserProducer  messaging.UserProducer
	NotifProducer messaging.NotifProducer

	// client
	S3Client storage.S3Client
}

func NewUserUsecase(
	Cfg *config.Config,
	DB *gorm.DB,

	// repository
	UserRepository repository.UserRepository,
	UserStatRepository repository.UserStatRepository,
	FollowRepository repository.FollowRepository,

	// producer
	UserProducer messaging.UserProducer,
	NotifProducer messaging.NotifProducer,

	// client
	S3Client storage.S3Client,
) *UserUsecaseImpl {
	return &UserUsecaseImpl{
		Config: Cfg,
		DB:     DB,

		// repository
		UserRepository:     UserRepository,
		UserStatRepository: UserStatRepository,
		FollowRepository:   FollowRepository,

		// producer
		UserProducer:  UserProducer,
		NotifProducer: NotifProducer,

		// client
		S3Client: S3Client,
	}
}
