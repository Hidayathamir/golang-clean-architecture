package config

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/image"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/notif"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/IBM/sarama"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Usecases struct {
	UserUsecase  user.UserUsecase
	ImageUsecase image.ImageUsecase
	NotifUsecase notif.NotifUsecase
}

func SetupUsecases(
	viperConfig *viper.Viper,
	db *gorm.DB,
	producer sarama.SyncProducer,
	awsS3Client *s3.Client,
) *Usecases {
	// setup repositories
	var userRepository repository.UserRepository
	userRepository = repository.NewUserRepository(viperConfig)
	userRepository = repository.NewUserRepositoryMwLogger(userRepository)

	var imageRepository repository.ImageRepository
	imageRepository = repository.NewImageRepository(viperConfig)
	imageRepository = repository.NewImageRepositoryMwLogger(imageRepository)

	var likeRepository repository.LikeRepository
	likeRepository = repository.NewLikeRepository(viperConfig)
	likeRepository = repository.NewLikeRepositoryMwLogger(likeRepository)

	var commentRepository repository.CommentRepository
	commentRepository = repository.NewCommentRepository(viperConfig)
	commentRepository = repository.NewCommentRepositoryMwLogger(commentRepository)

	var followRepository repository.FollowRepository
	followRepository = repository.NewFollowRepository(viperConfig)
	followRepository = repository.NewFollowRepositoryMwLogger(followRepository)

	// setup producer
	var userProducer messaging.UserProducer
	userProducer = messaging.NewUserProducer(viperConfig, producer)
	userProducer = messaging.NewUserProducerMwLogger(userProducer)

	var imageProducer messaging.ImageProducer
	imageProducer = messaging.NewImageProducer(viperConfig, producer)
	imageProducer = messaging.NewImageProducerMwLogger(imageProducer)

	var notifProducer messaging.NotifProducer
	notifProducer = messaging.NewNotifProducer(viperConfig, producer)
	notifProducer = messaging.NewNotifProducerMwLogger(notifProducer)

	// setup client
	var s3Client rest.S3Client
	s3Client = rest.NewS3Client(viperConfig, awsS3Client)
	s3Client = rest.NewS3ClientMwLogger(s3Client)

	// setup use cases
	var userUsecase user.UserUsecase
	userUsecase = user.NewUserUsecase(viperConfig, db, userRepository, followRepository, userProducer, notifProducer, s3Client)
	userUsecase = user.NewUserUsecaseMwLogger(userUsecase)

	var imageUsecase image.ImageUsecase
	imageUsecase = image.NewImageUsecase(viperConfig, db, imageRepository, likeRepository, commentRepository, followRepository, userRepository, imageProducer, notifProducer, s3Client)
	imageUsecase = image.NewImageUsecaseMwLogger(imageUsecase)

	var notifUsecase notif.NotifUsecase
	notifUsecase = notif.NewNotifUsecase(viperConfig, db)
	notifUsecase = notif.NewNotifUsecaseMwLogger(notifUsecase)

	return &Usecases{
		UserUsecase:  userUsecase,
		ImageUsecase: imageUsecase,
		NotifUsecase: notifUsecase,
	}
}

type Controllers struct {
	UserController  *http.UserController
	ImageController *http.ImageController
}

func SetupControllers(viperConfig *viper.Viper, usecases *Usecases) *Controllers {
	userController := http.NewUserController(viperConfig, usecases.UserUsecase)
	imageController := http.NewImageController(viperConfig, usecases.ImageUsecase)

	return &Controllers{
		UserController:  userController,
		ImageController: imageController,
	}
}

type Middlewares struct {
	AuthMiddleware      fiber.Handler
	TraceIDMiddleware   fiber.Handler
	OtelFiberMiddleware fiber.Handler
}

func SetupMiddlewares(usecases *Usecases) *Middlewares {
	authMiddleware := middleware.NewAuth(usecases.UserUsecase)
	traceIDMiddleware := middleware.NewTraceID()
	otelFiberMiddleware := middleware.NewOtelFiberMiddleware()

	return &Middlewares{
		AuthMiddleware:      authMiddleware,
		TraceIDMiddleware:   traceIDMiddleware,
		OtelFiberMiddleware: otelFiberMiddleware,
	}
}
