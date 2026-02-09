package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/storage"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseImage.go -pkg=mock . ImageUsecase

type ImageUsecase interface {
	Upload(ctx context.Context, req *model.UploadImageRequest) (*model.ImageResponse, error)
	Like(ctx context.Context, req *model.LikeImageRequest) error
	Comment(ctx context.Context, req *model.CommentImageRequest) error
	GetImage(ctx context.Context, req *model.GetImageRequest) (*model.ImageResponse, error)
	GetLike(ctx context.Context, req *model.GetLikeRequest) (*model.LikeResponseList, error)
	GetComment(ctx context.Context, req *model.GetCommentRequest) (*model.CommentResponseList, error)
	NotifyFollowerOnUpload(ctx context.Context, req *model.NotifyFollowerOnUploadRequest) error
	NotifyUserImageCommented(ctx context.Context, req *model.NotifyUserImageCommentedRequest) error
	BatchUpdateImageCommentCount(ctx context.Context, req *model.BatchUpdateImageCommentCountRequest) error
	NotifyUserImageLiked(ctx context.Context, req *model.NotifyUserImageLikedRequest) error
	BatchUpdateImageLikeCount(ctx context.Context, req *model.BatchUpdateImageLikeCountRequest) error
}

var _ ImageUsecase = &ImageUsecaseImpl{}

type ImageUsecaseImpl struct {
	Cfg *config.Config
	DB  *gorm.DB

	// repository
	ImageRepository   repository.ImageRepository
	LikeRepository    repository.LikeRepository
	CommentRepository repository.CommentRepository
	FollowRepository  repository.FollowRepository
	UserRepository    repository.UserRepository

	// producer
	ImageProducer messaging.ImageProducer
	NotifProducer messaging.NotifProducer

	// client
	S3Client storage.S3Client
}

func NewImageUsecase(
	Config *config.Config,
	DB *gorm.DB,

	// repository
	ImageRepository repository.ImageRepository,
	LikeRepository repository.LikeRepository,
	CommentRepository repository.CommentRepository,
	FollowRepository repository.FollowRepository,
	UserRepository repository.UserRepository,

	// producer
	ImageProducer messaging.ImageProducer,
	NotifProducer messaging.NotifProducer,

	// client
	S3Client storage.S3Client,
) *ImageUsecaseImpl {
	return &ImageUsecaseImpl{
		Cfg: Config,
		DB:  DB,

		// repository
		ImageRepository:   ImageRepository,
		LikeRepository:    LikeRepository,
		CommentRepository: CommentRepository,
		FollowRepository:  FollowRepository,
		UserRepository:    UserRepository,

		// producer
		ImageProducer: ImageProducer,
		NotifProducer: NotifProducer,

		// client
		S3Client: S3Client,
	}
}
