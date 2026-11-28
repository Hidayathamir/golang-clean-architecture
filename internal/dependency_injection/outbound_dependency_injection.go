package dependency_injection

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/cache"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/search"
	"github.com/Hidayathamir/golang-clean-architecture/internal/outbound/storage"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
)

type Caches struct {
	UserCache             cache.UserCache
	NotificationPublisher cache.NotificationPublisher
}

func SetupCaches(redisClient *redis.Client) *Caches {
	var userCache cache.UserCache
	userCache = cache.NewUserCache(redisClient)
	userCache = cache.NewUserCacheMwLogger(userCache)

	var notificationPublisher cache.NotificationPublisher
	notificationPublisher = cache.NewNotificationPublisher(redisClient)
	notificationPublisher = cache.NewNotificationPublisherMwLogger(notificationPublisher)

	return &Caches{
		UserCache:             userCache,
		NotificationPublisher: notificationPublisher,
	}
}

type Messagings struct {
	UserProducer  messaging.UserProducer
	ImageProducer messaging.ImageProducer
	NotifProducer messaging.NotifProducer
}

func SetupMessagings(cfg *config.Config, repositories *Repositories) *Messagings {
	var userProducer messaging.UserProducer
	userProducer = messaging.NewUserProducer(cfg, repositories.OutboxRepository)
	userProducer = messaging.NewUserProducerMwLogger(userProducer)

	var imageProducer messaging.ImageProducer
	imageProducer = messaging.NewImageProducer(cfg, repositories.OutboxRepository)
	imageProducer = messaging.NewImageProducerMwLogger(imageProducer)

	var notifProducer messaging.NotifProducer
	notifProducer = messaging.NewNotifProducer(cfg, repositories.OutboxRepository)
	notifProducer = messaging.NewNotifProducerMwLogger(notifProducer)

	return &Messagings{
		UserProducer:  userProducer,
		ImageProducer: imageProducer,
		NotifProducer: notifProducer,
	}
}

type Repositories struct {
	UserRepository         repository.UserRepository
	NotificationRepository repository.NotificationRepository
	UserStatRepository     repository.UserStatRepository
	ImageRepository        repository.ImageRepository
	LikeRepository         repository.LikeRepository
	CommentRepository      repository.CommentRepository
	FollowRepository       repository.FollowRepository
	OutboxRepository       repository.OutboxRepository
	IdempotencyRepository  repository.IdempotencyRepository
	ChatRepository         repository.ChatRepository
}

func SetupRepositories(cfg *config.Config) *Repositories {
	var userRepository repository.UserRepository
	userRepository = repository.NewUserRepository(cfg)
	userRepository = repository.NewUserRepositoryMwLogger(userRepository)

	var notificationRepository repository.NotificationRepository
	notificationRepository = repository.NewNotificationRepository(cfg)
	notificationRepository = repository.NewNotificationRepositoryMwLogger(notificationRepository)

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

	var idempotencyRepository repository.IdempotencyRepository
	idempotencyRepository = repository.NewIdempotencyRepository(cfg)
	idempotencyRepository = repository.NewIdempotencyRepositoryMwLogger(idempotencyRepository)

	var chatRepository repository.ChatRepository
	chatRepository = repository.NewChatRepository(cfg)
	chatRepository = repository.NewChatRepositoryMwLogger(chatRepository)

	return &Repositories{
		UserRepository:         userRepository,
		NotificationRepository: notificationRepository,
		UserStatRepository:     userStatRepository,
		ImageRepository:        imageRepository,
		LikeRepository:         likeRepository,
		CommentRepository:      commentRepository,
		FollowRepository:       followRepository,
		OutboxRepository:       outboxRepository,
		IdempotencyRepository:  idempotencyRepository,
		ChatRepository:         chatRepository,
	}
}

type Searches struct {
	ImageSearch search.ImageSearch
}

func SetupSearches(elasticsearchClient *elasticsearch.Client) *Searches {
	var imageSearch search.ImageSearch
	imageSearch = search.NewImageSearch(elasticsearchClient)
	imageSearch = search.NewImageSearchMwLogger(imageSearch)

	return &Searches{
		ImageSearch: imageSearch,
	}
}

type Storages struct {
	S3Client storage.S3Client
}

func SetupStorages(cfg *config.Config, awsS3Client *s3.Client) *Storages {
	var s3Client storage.S3Client
	s3Client = storage.NewS3Client(cfg, awsS3Client)
	s3Client = storage.NewS3ClientMwLogger(s3Client)

	return &Storages{
		S3Client: s3Client,
	}
}
