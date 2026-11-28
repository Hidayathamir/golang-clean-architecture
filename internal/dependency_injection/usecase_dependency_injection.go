package dependency_injection

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/inbound/websocket"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/cache"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/search"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/storage"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/idempotencyusecase"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/imageusecase"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/notifusecase"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/userusecase"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
	"gorm.io/gorm"
)

type Usecases struct {
	UserUsecase        userusecase.UserUsecase
	ImageUsecase       imageusecase.ImageUsecase
	NotifUsecase       notifusecase.NotifUsecase
	NotificationHub    *websocket.Hub
	RedisClient        *redis.Client
	IdempotencyRepo    repository.IdempotencyRepository
	IdempotencyUsecase idempotencyusecase.IdempotencyUsecase
}

func SetupUsecases(
	cfg *config.Config,
	db *gorm.DB,
	producer *kgo.Client,
	awsS3Client *s3.Client,
	redisClient *redis.Client,
	elasticsearchClient *elasticsearch.Client,
) *Usecases {
	// setup repositories
	var userRepository repository.UserRepository
	userRepository = repository.NewUserRepository(cfg)
	userRepository = repository.NewUserRepositoryMwLogger(userRepository)

	var notificationRepository repository.NotificationRepository
	notificationRepository = repository.NewNotificationRepository(cfg)
	notificationRepository = repository.NewNotificationRepositoryMwLogger(notificationRepository)

	notificationHub := websocket.NewHub()

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

	var outboxRepository repository.OutboxRepository
	outboxRepository = repository.NewOutboxRepository(cfg)
	outboxRepository = repository.NewOutboxRepositoryMwLogger(outboxRepository)

	var idempotencyRepo repository.IdempotencyRepository
	idempotencyRepo = repository.NewIdempotencyRepository(cfg)
	idempotencyRepo = repository.NewIdempotencyRepositoryMwLogger(idempotencyRepo)

	// setup idempotency usecase
	var idempotencyUsecase idempotencyusecase.IdempotencyUsecase
	idempotencyUsecase = idempotencyusecase.NewIdempotencyUsecase(cfg, db, idempotencyRepo)
	idempotencyUsecase = idempotencyusecase.NewIdempotencyUsecaseMwLogger(idempotencyUsecase)

	// setup producer
	var userProducer messaging.UserProducer
	userProducer = messaging.NewUserProducer(cfg, outboxRepository)
	userProducer = messaging.NewUserProducerMwLogger(userProducer)

	var imageProducer messaging.ImageProducer
	imageProducer = messaging.NewImageProducer(cfg, outboxRepository)
	imageProducer = messaging.NewImageProducerMwLogger(imageProducer)

	var notifProducer messaging.NotifProducer
	notifProducer = messaging.NewNotifProducer(cfg, outboxRepository)
	notifProducer = messaging.NewNotifProducerMwLogger(notifProducer)

	// setup client
	var s3Client storage.S3Client
	s3Client = storage.NewS3Client(cfg, awsS3Client)
	s3Client = storage.NewS3ClientMwLogger(s3Client)

	// setup cache
	var userCache cache.UserCache
	userCache = cache.NewUserCache(redisClient)
	userCache = cache.NewUserCacheMwLogger(userCache)

	var notificationPublisher cache.NotificationPublisher
	notificationPublisher = cache.NewNotificationPublisher(redisClient)
	notificationPublisher = cache.NewNotificationPublisherMwLogger(notificationPublisher)

	// setup search
	var imageSearch search.ImageSearch
	imageSearch = search.NewImageSearch(elasticsearchClient)
	imageSearch = search.NewImageSearchMwLogger(imageSearch)

	// setup use cases
	var notifUsecase notifusecase.NotifUsecase
	notifUsecase = notifusecase.NewNotifUsecase(cfg, db, notificationRepository, userRepository, notificationPublisher)
	notifUsecase = notifusecase.NewNotifUsecaseMwLogger(notifUsecase)

	var userUsecase userusecase.UserUsecase
	userUsecase = userusecase.NewUserUsecase(cfg, db, userRepository, userStatRepository, followRepository, userProducer, notifProducer, s3Client, userCache, notifUsecase)
	userUsecase = userusecase.NewUserUsecaseMwLogger(userUsecase)

	var imageUsecase imageusecase.ImageUsecase
	imageUsecase = imageusecase.NewImageUsecase(cfg, db, imageRepository, likeRepository, commentRepository, followRepository, userRepository, imageProducer, notifProducer, s3Client, imageSearch, notifUsecase)
	imageUsecase = imageusecase.NewImageUsecaseMwLogger(imageUsecase)

	return &Usecases{
		UserUsecase:        userUsecase,
		ImageUsecase:       imageUsecase,
		NotifUsecase:       notifUsecase,
		NotificationHub:    notificationHub,
		RedisClient:        redisClient,
		IdempotencyRepo:    idempotencyRepo,
		IdempotencyUsecase: idempotencyUsecase,
	}
}
