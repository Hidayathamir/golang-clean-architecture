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
	db *gorm.DB,
	app *fiber.App,
	validate *validator.Validate,
	producer sarama.SyncProducer,
) *Usecases {
	// setup repositories
	var userRepository repository.UserRepository
	userRepository = repository.NewUserRepository(viperConfig)
	userRepository = repository.NewUserRepositoryMwLogger(userRepository)

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
	s3Client = rest.NewS3Client(viperConfig)
	s3Client = rest.NewS3ClientMwLogger(s3Client)

	var slackClient rest.SlackClient
	slackClient = rest.NewSlackClient(viperConfig)
	slackClient = rest.NewSlackClientMwLogger(slackClient)

	// setup use cases
	var userUsecase user.UserUsecase
	userUsecase = user.NewUserUsecase(viperConfig, db, validate, userRepository, userProducer, s3Client, slackClient)
	userUsecase = user.NewUserUsecaseMwLogger(userUsecase)

	var contactUsecase contact.ContactUsecase
	contactUsecase = contact.NewContactUsecase(viperConfig, db, validate, contactRepository, contactProducer, slackClient)
	contactUsecase = contact.NewContactUsecaseMwLogger(contactUsecase)

	var addressUsecase address.AddressUsecase
	addressUsecase = address.NewAddressUsecase(viperConfig, db, validate, contactRepository, addressRepository, addressProducer, paymentClient)
	addressUsecase = address.NewAddressUsecaseMwLogger(addressUsecase)

	var todoUsecase todo.TodoUsecase
	todoUsecase = todo.NewTodoUsecase(viperConfig, db, validate, todoRepository, todoProducer)
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

func SetupControllers(viperConfig *viper.Viper, usecases *Usecases) *Controllers {
	userController := http.NewUserController(viperConfig, usecases.UserUsecase)
	contactController := http.NewContactController(viperConfig, usecases.ContactUsecase)
	addressController := http.NewAddressController(viperConfig, usecases.AddressUsecase)
	todoController := http.NewTodoController(viperConfig, usecases.TodoUsecase)

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
