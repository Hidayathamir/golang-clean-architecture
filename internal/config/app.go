package config

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/contact"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/todo"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user"
	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Usecases struct {
	UserUsecase    user.UserUsecase
	ContactUsecase contact.ContactUsecase
	AddressUsecase address.AddressUsecase
	TodoUsecase    todo.TodoUsecase
}

func SetupUsecases(
	viperConfig *viper.Viper,
	log *logrus.Logger,
	db *gorm.DB,
	app *fiber.App,
	validate *validator.Validate,
	producer sarama.SyncProducer,
) *Usecases {
	// setup repositories
	var userRepository repository.UserRepository
	userRepository = repository.NewUserRepository(viperConfig, log)
	userRepository = repository.NewUserRepositoryMwLogger(userRepository)

	var contactRepository repository.ContactRepository
	contactRepository = repository.NewContactRepository(viperConfig, log)
	contactRepository = repository.NewContactRepositoryMwLogger(contactRepository)

	var addressRepository repository.AddressRepository
	addressRepository = repository.NewAddressRepository(viperConfig, log)
	addressRepository = repository.NewAddressRepositoryMwLogger(addressRepository)

	var todoRepository repository.TodoRepository
	todoRepository = repository.NewTodoRepository(viperConfig, log)
	todoRepository = repository.NewTodoRepositoryMwLogger(todoRepository)

	// setup producer
	var userProducer messaging.UserProducer
	userProducer = messaging.NewUserProducer(viperConfig, log, producer)
	userProducer = messaging.NewUserProducerMwTelemetry(userProducer)
	userProducer = messaging.NewUserProducerMwLogger(userProducer)

	var contactProducer messaging.ContactProducer
	contactProducer = messaging.NewContactProducer(viperConfig, log, producer)
	contactProducer = messaging.NewContactProducerMwTelemetry(contactProducer)
	contactProducer = messaging.NewContactProducerMwLogger(contactProducer)

	var addressProducer messaging.AddressProducer
	addressProducer = messaging.NewAddressProducer(viperConfig, log, producer)
	addressProducer = messaging.NewAddressProducerMwTelemetry(addressProducer)
	addressProducer = messaging.NewAddressProducerMwLogger(addressProducer)

	var todoProducer messaging.TodoProducer
	todoProducer = messaging.NewTodoProducer(viperConfig, log, producer)
	todoProducer = messaging.NewTodoProducerMwTelemetry(todoProducer)
	todoProducer = messaging.NewTodoProducerMwLogger(todoProducer)

	// setup client
	var paymentClient rest.PaymentClient
	paymentClient = rest.NewPaymentClient(viperConfig)
	paymentClient = rest.NewPaymentClientMwLogger(paymentClient)

	var s3Client rest.S3Client
	s3Client = rest.NewS3Client(viperConfig)
	s3Client = rest.NewS3ClientMwLogger(s3Client)

	var slackClient rest.SlackClient
	slackClient = rest.NewSlackClient(viperConfig)
	slackClient = rest.NewSlackClientMwLogger(slackClient)

	// setup use cases
	var userUsecase user.UserUsecase
	userUsecase = user.NewUserUsecase(viperConfig, log, db, validate, userRepository, userProducer, s3Client, slackClient)
	userUsecase = user.NewUserUsecaseMwLogger(userUsecase)

	var contactUsecase contact.ContactUsecase
	contactUsecase = contact.NewContactUsecase(viperConfig, log, db, validate, contactRepository, contactProducer, slackClient)
	contactUsecase = contact.NewContactUsecaseMwLogger(contactUsecase)

	var addressUsecase address.AddressUsecase
	addressUsecase = address.NewAddressUsecase(viperConfig, log, db, validate, contactRepository, addressRepository, addressProducer, paymentClient)
	addressUsecase = address.NewAddressUsecaseMwLogger(addressUsecase)

	var todoUsecase todo.TodoUsecase
	todoUsecase = todo.NewTodoUsecase(viperConfig, log, db, validate, todoRepository, todoProducer)
	todoUsecase = todo.NewTodoUsecaseMwLogger(todoUsecase)

	return &Usecases{
		UserUsecase:    userUsecase,
		ContactUsecase: contactUsecase,
		AddressUsecase: addressUsecase,
		TodoUsecase:    todoUsecase,
	}
}

type Controllers struct {
	UserController    *http.UserController
	ContactController *http.ContactController
	AddressController *http.AddressController
	TodoController    *http.TodoController
}

func SetupControllers(viperConfig *viper.Viper, log *logrus.Logger, usecases *Usecases) *Controllers {
	userController := http.NewUserController(viperConfig, log, usecases.UserUsecase)
	contactController := http.NewContactController(viperConfig, log, usecases.ContactUsecase)
	addressController := http.NewAddressController(viperConfig, log, usecases.AddressUsecase)
	todoController := http.NewTodoController(viperConfig, log, usecases.TodoUsecase)

	return &Controllers{
		UserController:    userController,
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
