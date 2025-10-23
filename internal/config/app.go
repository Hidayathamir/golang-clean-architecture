package config

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http/middleware"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/messaging"
	"github.com/Hidayathamir/golang-clean-architecture/internal/gateway/rest"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/address"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/contact"
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
}

func SetupUsecases(
	db *gorm.DB,
	app *fiber.App,
	log *logrus.Logger,
	validate *validator.Validate,
	viperConfig *viper.Viper,
	producer sarama.SyncProducer,
) *Usecases {
	// setup repositories
	var userRepository repository.UserRepository
	userRepository = repository.NewUserRepository(log)
	userRepository = repository.NewUserRepositoryMwLogger(userRepository)

	var contactRepository repository.ContactRepository
	contactRepository = repository.NewContactRepository(log)
	contactRepository = repository.NewContactRepositoryMwLogger(contactRepository)

	var addressRepository repository.AddressRepository
	addressRepository = repository.NewAddressRepository(log)
	addressRepository = repository.NewAddressRepositoryMwLogger(addressRepository)

	// setup producer
	var userProducer messaging.UserProducer
	userProducer = messaging.NewUserProducer(producer, log)
	userProducer = messaging.NewUserProducerMwLogger(userProducer)

	var contactProducer messaging.ContactProducer
	contactProducer = messaging.NewContactProducer(producer, log)
	contactProducer = messaging.NewContactProducerMwLogger(contactProducer)

	var addressProducer messaging.AddressProducer
	addressProducer = messaging.NewAddressProducer(producer, log)
	addressProducer = messaging.NewAddressProducerMwLogger(addressProducer)

	// setup client
	var paymentClient rest.PaymentClient
	paymentClient = rest.NewPaymentClient()
	paymentClient = rest.NewPaymentClientMwLogger(paymentClient)

	var s3Client rest.S3Client
	s3Client = rest.NewS3Client()
	s3Client = rest.NewS3ClientMwLogger(s3Client)

	var slackClient rest.SlackClient
	slackClient = rest.NewSlackClient()
	slackClient = rest.NewSlackClientMwLogger(slackClient)

	// setup use cases
	var userUsecase user.UserUsecase
	userUsecase = user.NewUserUsecase(db, log, validate, userRepository, userProducer, s3Client, slackClient)
	userUsecase = user.NewUserUsecaseMwLogger(userUsecase)

	var contactUsecase contact.ContactUsecase
	contactUsecase = contact.NewContactUsecase(db, log, validate, contactRepository, contactProducer, slackClient)
	contactUsecase = contact.NewContactUsecaseMwLogger(contactUsecase)

	var addressUsecase address.AddressUsecase
	addressUsecase = address.NewAddressUsecase(db, log, validate, contactRepository, addressRepository, addressProducer, paymentClient)
	addressUsecase = address.NewAddressUsecaseMwLogger(addressUsecase)

	return &Usecases{
		UserUsecase:    userUsecase,
		ContactUsecase: contactUsecase,
		AddressUsecase: addressUsecase,
	}
}

type Controllers struct {
	UserController    *http.UserController
	ContactController *http.ContactController
	AddressController *http.AddressController
}

func SetupControllers(usecases *Usecases, log *logrus.Logger) *Controllers {
	userController := http.NewUserController(usecases.UserUsecase, log)
	contactController := http.NewContactController(usecases.ContactUsecase, log)
	addressController := http.NewAddressController(usecases.AddressUsecase, log)

	return &Controllers{
		UserController:    userController,
		ContactController: contactController,
		AddressController: addressController,
	}
}

type Middlewares struct {
	AuthMiddleware    fiber.Handler
	TraceIDMiddleware fiber.Handler
}

func SetupMiddlewares(usecases *Usecases) *Middlewares {
	authMiddleware := middleware.NewAuth(usecases.UserUsecase)
	traceIDMiddleware := middleware.NewTraceID()

	return &Middlewares{
		AuthMiddleware:    authMiddleware,
		TraceIDMiddleware: traceIDMiddleware,
	}
}
