package dependency_injection

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/cache"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/search"
	"github.com/Hidayathamir/golang-clean-architecture/internal/infra/storage"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/image"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/notif"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
	"gorm.io/gorm"
)

type Usecases struct {
	UserUsecase  user.UserUsecase
	ImageUsecase image.ImageUsecase
	NotifUsecase notif.NotifUsecase
}

func SetupUsecases(
	cfg *config.Config,
	db *gorm.DB,
	producer *kgo.Client,
	awsS3Client *s3.Client,
	redisClient *redis.Client,
	opensearchClient *opensearch.Client,
) *Usecases {
	// setup repositories
	var userRepository repository.UserRepository
	userRepository = repository.NewUserRepository(cfg)
	userRepository = repository.NewUserRepositoryMwLogger(userRepository)

	var userStatRepository repository.UserStatRepository
	userStatRepository = repository.NewUserStatRepository(cfg)
	userStatRepository = repository.NewUserStatRepositoryMwLogger(userStatRepository)

	var imageRepository repository.ImageRepository
	imageRepository = repository.NewImageRepository(cfg)
	imageRepository = repository.NewImageRepositoryMwLogger(imageRepository)

	var likeRepository repository.LikeRepository
	likeRepository = repository.NewLikeRepository(cfg)
	likeRepository = repository.NewLikeRepositoryMwLogger(likeRepository)

	var commentRepository repository.CommentRepository
	commentRepository = repository.NewCommentRepository(cfg)
	commentRepository = repository.NewCommentRepositoryMwLogger(commentRepository)

	var followRepository repository.FollowRepository
	followRepository = repository.NewFollowRepository(cfg)
	followRepository = repository.NewFollowRepositoryMwLogger(followRepository)

	// setup producer
	var userProducer messaging.UserProducer
	userProducer = messaging.NewUserProducer(cfg, producer)
	userProducer = messaging.NewUserProducerMwLogger(userProducer)

	var imageProducer messaging.ImageProducer
	imageProducer = messaging.NewImageProducer(cfg, producer)
	imageProducer = messaging.NewImageProducerMwLogger(imageProducer)

	var notifProducer messaging.NotifProducer
	notifProducer = messaging.NewNotifProducer(cfg, producer)
	notifProducer = messaging.NewNotifProducerMwLogger(notifProducer)

	// setup client
	var s3Client storage.S3Client
	s3Client = storage.NewS3Client(cfg, awsS3Client)
	s3Client = storage.NewS3ClientMwLogger(s3Client)

	// setup cache
	var userCache cache.UserCache
	userCache = cache.NewUserCache(redisClient)
	userCache = cache.NewUserCacheMwLogger(userCache)

	// setup search
	var imageSearch search.ImageSearch2
	imageSearch = search.NewImageSearch2(opensearchClient)
	imageSearch = search.NewImageSearch2MwLogger(imageSearch)

	// setup use cases
	var userUsecase user.UserUsecase
	userUsecase = user.NewUserUsecase(cfg, db, userRepository, userStatRepository, followRepository, userProducer, notifProducer, s3Client, userCache)
	userUsecase = user.NewUserUsecaseMwLogger(userUsecase)

	var imageUsecase image.ImageUsecase
	imageUsecase = image.NewImageUsecase(cfg, db, imageRepository, likeRepository, commentRepository, followRepository, userRepository, imageProducer, notifProducer, s3Client, imageSearch)
	imageUsecase = image.NewImageUsecaseMwLogger(imageUsecase)

	var notifUsecase notif.NotifUsecase
	notifUsecase = notif.NewNotifUsecase(cfg, db)
	notifUsecase = notif.NewNotifUsecaseMwLogger(notifUsecase)

	return &Usecases{
		UserUsecase:  userUsecase,
		ImageUsecase: imageUsecase,
		NotifUsecase: notifUsecase,
	}
}
