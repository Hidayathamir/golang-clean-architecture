package imageusecase

import (
	"context"

	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/dto"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/search"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/storage"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/MockUsecaseImage.go -pkg=mock . ImageUsecase

type ImageUsecase interface {
	PresignUpload(ctx context.Context, req dto.PresignUploadImageRequest) (dto.PresignUploadImageResponse, error)
	ConfirmUpload(ctx context.Context, req dto.ConfirmUploadImageRequest) (dto.ImageResponse, error)
	Like(ctx context.Context, req dto.LikeImageRequest) error
	Comment(ctx context.Context, req dto.CommentImageRequest) error
	GetImage(ctx context.Context, req dto.GetImageRequest) (dto.ImageResponse, error)
	GetFeed(ctx context.Context, req dto.FeedRequest) (dto.FeedResponse, error)
	GetUserImages(ctx context.Context, req dto.UserImagesRequest) (dto.UserImagesResponse, error)
	Explore(ctx context.Context, req dto.ExploreRequest) (dto.ExploreResponse, error)
	DeleteImage(ctx context.Context, req dto.DeleteImageRequest) error
	UpdateImage(ctx context.Context, req dto.UpdateImageRequest) error
	Unlike(ctx context.Context, req dto.UnlikeImageRequest) error
	DeleteImageComment(ctx context.Context, req dto.DeleteImageCommentRequest) error
	GetLike(ctx context.Context, req dto.GetLikeRequest) (dto.LikesResponse, error)
	GetComment(ctx context.Context, req dto.GetCommentRequest) (dto.CommentsResponse, error)
	NotifyFollowerOnUpload(ctx context.Context, req dto.NotifyFollowerOnUploadRequest) error
	SyncImageToElasticsearch(ctx context.Context, req dto.SyncImageToElasticsearchRequest) error
	NotifyUserImageCommented(ctx context.Context, req dto.NotifyUserImageCommentedRequest) error
	BatchUpdateImageCommentCount(ctx context.Context, req dto.BatchUpdateImageCommentCountRequest) error
	NotifyUserImageLiked(ctx context.Context, req dto.NotifyUserImageLikedRequest) error
	BatchUpdateImageLikeCount(ctx context.Context, req dto.BatchUpdateImageLikeCountRequest) error
	CleanupDeletedImageStorage(ctx context.Context, req dto.CleanupDeletedImageStorageRequest) error
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

	// storage
	S3Client storage.S3Client

	// search
	ImageSearch search.ImageSearch
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

	// storage
	S3Client storage.S3Client,

	// search
	ImageSearch search.ImageSearch,
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

		// storage
		S3Client: S3Client,

		// search
		ImageSearch: ImageSearch,
	}
}
