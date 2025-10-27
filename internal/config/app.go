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
	user2usecase "github.com/Hidayathamir/golang-clean-architecture/internal/usecase/user2"
	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Usecases struct {
	UserUsecase    user.UserUsecase
	User2Usecase   user2usecase.User2Usecase
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

	var user2Repository repository.User2Repository
	user2Repository = repository.NewUser2Repository(viperConfig, log)
	user2Repository = repository.NewUser2RepositoryMwLogger(user2Repository)

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
	userProducer = messaging.NewUserProducerMwLogger(userProducer)

	var contactProducer messaging.ContactProducer
	contactProducer = messaging.NewContactProducer(viperConfig, log, producer)
	contactProducer = messaging.NewContactProducerMwLogger(contactProducer)

	var addressProducer messaging.AddressProducer
	addressProducer = messaging.NewAddressProducer(viperConfig, log, producer)
	addressProducer = messaging.NewAddressProducerMwLogger(addressProducer)

	var todoProducer messaging.TodoProducer
	todoProducer = messaging.NewTodoProducer(viperConfig, log, producer)
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

	var user2Usecase user2usecase.User2Usecase = user2usecase.NewUser2Usecase(viperConfig, log, db, validate, user2Repository)

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
		User2Usecase:   user2Usecase,
		ContactUsecase: contactUsecase,
		AddressUsecase: addressUsecase,
		TodoUsecase:    todoUsecase,
	}
}

type Controllers struct {
	UserController    *http.UserController
	User2Controller   *http.User2Controller
	ContactController *http.ContactController
	AddressController *http.AddressController
	TodoController    *http.TodoController
}

func SetupControllers(viperConfig *viper.Viper, log *logrus.Logger, usecases *Usecases) *Controllers {
	userController := http.NewUserController(viperConfig, log, usecases.UserUsecase)
	user2Controller := http.NewUser2Controller(viperConfig, log, usecases.User2Usecase)
	contactController := http.NewContactController(viperConfig, log, usecases.ContactUsecase)
	addressController := http.NewAddressController(viperConfig, log, usecases.AddressUsecase)
	todoController := http.NewTodoController(viperConfig, log, usecases.TodoUsecase)

	return &Controllers{
		UserController:    userController,
		User2Controller:   user2Controller,
		ContactController: contactController,
		AddressController: addressController,
		TodoController:    todoController,
	}
}

type Middlewares struct {
	AuthMiddleware      fiber.Handler
	User2AuthMiddleware fiber.Handler
	TraceIDMiddleware   fiber.Handler
}

func SetupMiddlewares(usecases *Usecases) *Middlewares {
	authMiddleware := middleware.NewAuth(usecases.UserUsecase)
	user2AuthMiddleware := middleware.NewUser2Auth(usecases.User2Usecase)
	traceIDMiddleware := middleware.NewTraceID()

	return &Middlewares{
		AuthMiddleware:      authMiddleware,
		User2AuthMiddleware: user2AuthMiddleware,
		TraceIDMiddleware:   traceIDMiddleware,
	}
}
