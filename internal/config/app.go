package config

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/contact"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/image"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/notif"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/todo"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/IBM/sarama"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Usecases struct {
	UserUsecase    user.UserUsecase
	ImageUsecase   image.ImageUsecase
	NotifUsecase   notif.NotifUsecase
	ContactUsecase contact.ContactUsecase
	AddressUsecase address.AddressUsecase
	TodoUsecase    todo.TodoUsecase
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

	var contactRepository repository.ContactRepository
	contactRepository = repository.NewContactRepository(viperConfig)
	contactRepository = repository.NewContactRepositoryMwLogger(contactRepository)

	var addressRepository repository.AddressRepository
	addressRepository = repository.NewAddressRepository(viperConfig)
	addressRepository = repository.NewAddressRepositoryMwLogger(addressRepository)

	var todoRepository repository.TodoRepository
	todoRepository = repository.NewTodoRepository(viperConfig)
	todoRepository = repository.NewTodoRepositoryMwLogger(todoRepository)

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

	var contactProducer messaging.ContactProducer
	contactProducer = messaging.NewContactProducer(viperConfig, producer)
	contactProducer = messaging.NewContactProducerMwLogger(contactProducer)

	var addressProducer messaging.AddressProducer
	addressProducer = messaging.NewAddressProducer(viperConfig, producer)
	addressProducer = messaging.NewAddressProducerMwLogger(addressProducer)

	var todoProducer messaging.TodoProducer
	todoProducer = messaging.NewTodoProducer(viperConfig, producer)
	todoProducer = messaging.NewTodoProducerMwLogger(todoProducer)

	// setup client
	var paymentClient rest.PaymentClient
	paymentClient = rest.NewPaymentClient(viperConfig)
	paymentClient = rest.NewPaymentClientMwLogger(paymentClient)

	var s3Client rest.S3Client
	s3Client = rest.NewS3Client(viperConfig, awsS3Client)
	s3Client = rest.NewS3ClientMwLogger(s3Client)

	var slackClient rest.SlackClient
	slackClient = rest.NewSlackClient(viperConfig)
	slackClient = rest.NewSlackClientMwLogger(slackClient)

	// setup use cases
	var userUsecase user.UserUsecase
	userUsecase = user.NewUserUsecase(viperConfig, db, userRepository, followRepository, userProducer, s3Client, slackClient)
	userUsecase = user.NewUserUsecaseMwLogger(userUsecase)

	var imageUsecase image.ImageUsecase
	imageUsecase = image.NewImageUsecase(viperConfig, db, imageRepository, likeRepository, commentRepository, followRepository, userRepository, imageProducer, notifProducer, s3Client)
	imageUsecase = image.NewImageUsecaseMwLogger(imageUsecase)

	var notifUsecase notif.NotifUsecase
	notifUsecase = notif.NewNotifUsecase(viperConfig, db)
	notifUsecase = notif.NewNotifUsecaseMwLogger(notifUsecase)

	var contactUsecase contact.ContactUsecase
	contactUsecase = contact.NewContactUsecase(viperConfig, db, contactRepository, contactProducer, slackClient)
	contactUsecase = contact.NewContactUsecaseMwLogger(contactUsecase)

	var addressUsecase address.AddressUsecase
	addressUsecase = address.NewAddressUsecase(viperConfig, db, contactRepository, addressRepository, addressProducer, paymentClient)
	addressUsecase = address.NewAddressUsecaseMwLogger(addressUsecase)

	var todoUsecase todo.TodoUsecase
	todoUsecase = todo.NewTodoUsecase(viperConfig, db, todoRepository, todoProducer)
	todoUsecase = todo.NewTodoUsecaseMwLogger(todoUsecase)

	return &Usecases{
		UserUsecase:    userUsecase,
		ImageUsecase:   imageUsecase,
		NotifUsecase:   notifUsecase,
		ContactUsecase: contactUsecase,
		AddressUsecase: addressUsecase,
		TodoUsecase:    todoUsecase,
	}
}

type Controllers struct {
	UserController    *http.UserController
	ImageController   *http.ImageController
	ContactController *http.ContactController
	AddressController *http.AddressController
	TodoController    *http.TodoController
}

func SetupControllers(viperConfig *viper.Viper, usecases *Usecases) *Controllers {
	userController := http.NewUserController(viperConfig, usecases.UserUsecase)
	imageController := http.NewImageController(viperConfig, usecases.ImageUsecase)
	contactController := http.NewContactController(viperConfig, usecases.ContactUsecase)
	addressController := http.NewAddressController(viperConfig, usecases.AddressUsecase)
	todoController := http.NewTodoController(viperConfig, usecases.TodoUsecase)

	return &Controllers{
		UserController:    userController,
		ImageController:   imageController,
		ContactController: contactController,
		AddressController: addressController,
		TodoController:    todoController,
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
