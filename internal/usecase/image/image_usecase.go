package image

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/ImageUsecase.go -pkg=mock . ImageUsecase

type ImageUsecase interface {
	Upload(ctx context.Context, req *model.UploadImageRequest) error
	Like(ctx context.Context, req *model.LikeImageRequest) error
	Comment(ctx context.Context, req *model.CommentImageRequest) error
}

var _ ImageUsecase = &ImageUsecaseImpl{}

type ImageUsecaseImpl struct {
	Config *viper.Viper
	DB     *gorm.DB

	// repository
	ImageRepository   repository.ImageRepository
	LikeRepository    repository.LikeRepository
	CommentRepository repository.CommentRepository

	// producer
	ImageUploadedProducer  messaging.ImageUploadedProducer
	ImageLikedProducer     messaging.ImageLikedProducer
	ImageCommentedProducer messaging.ImageCommentedProducer

	// client
	S3Client rest.S3Client
}

func NewImageUsecase(
	Config *viper.Viper,
	DB *gorm.DB,

	// repository
	ImageRepository repository.ImageRepository,
	LikeRepository repository.LikeRepository,
	CommentRepository repository.CommentRepository,

	// producer
	ImageUploadedProducer messaging.ImageUploadedProducer,
	ImageLikedProducer messaging.ImageLikedProducer,
	ImageCommentedProducer messaging.ImageCommentedProducer,

	// client
	S3Client rest.S3Client,
) *ImageUsecaseImpl {
	return &ImageUsecaseImpl{
		Config: Config,
		DB:     DB,

		// repository
		ImageRepository:   ImageRepository,
		LikeRepository:    LikeRepository,
		CommentRepository: CommentRepository,

		// producer
		ImageUploadedProducer:  ImageUploadedProducer,
		ImageLikedProducer:     ImageLikedProducer,
		ImageCommentedProducer: ImageCommentedProducer,

		// client
		S3Client: S3Client,
	}
}
