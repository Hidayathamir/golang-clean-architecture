package config

import (
	"golang-clean-architecture/internal/delivery/http"
	"golang-clean-architecture/internal/delivery/http/middleware"
	"golang-clean-architecture/internal/delivery/http/route"
	"golang-clean-architecture/internal/gateway/messaging"
	"golang-clean-architecture/internal/gateway/rest"
	"golang-clean-architecture/internal/repository"
	"golang-clean-architecture/internal/usecase/address"
	"golang-clean-architecture/internal/usecase/contact"
	"golang-clean-architecture/internal/usecase/user"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Producer sarama.SyncProducer
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	var userRepository repository.UserRepository
	userRepository = repository.NewUserRepository(config.Log)
	userRepository = repository.NewUserRepositoryMwLogger(config.Log, userRepository)

	var contactRepository repository.ContactRepository
	contactRepository = repository.NewContactRepository(config.Log)
	contactRepository = repository.NewContactRepositoryMwLogger(config.Log, contactRepository)

	var addressRepository repository.AddressRepository
	addressRepository = repository.NewAddressRepository(config.Log)
	addressRepository = repository.NewAddressRepositoryMwLogger(config.Log, addressRepository)

	// setup producer
	var userProducer messaging.UserProducer
	userProducer = messaging.NewUserProducer(config.Producer, config.Log)
	userProducer = messaging.NewUserProducerMwLogger(config.Log, userProducer)

	var contactProducer messaging.ContactProducer
	contactProducer = messaging.NewContactProducer(config.Producer, config.Log)
	contactProducer = messaging.NewContactProducerMwLogger(config.Log, contactProducer)

	var addressProducer messaging.AddressProducer
	addressProducer = messaging.NewAddressProducer(config.Producer, config.Log)
	addressProducer = messaging.NewAddressProducerMwLogger(config.Log, addressProducer)

	// setup client
	var paymentClient rest.PaymentClient
	paymentClient = rest.NewPaymentClient()
	paymentClient = rest.NewPaymentClientMwLogger(config.Log, paymentClient)

	var s3Client rest.S3Client
	s3Client = rest.NewS3Client()
	s3Client = rest.NewS3ClientMwLogger(config.Log, s3Client)

	var slackClient rest.SlackClient
	slackClient = rest.NewSlackClient()
	slackClient = rest.NewSlackClientMwLogger(config.Log, slackClient)

	// setup use cases
	var userUseCase user.UserUseCase
	userUseCase = user.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository, userProducer, s3Client, slackClient)
	userUseCase = user.NewUserUseCaseMwLogger(config.Log, userUseCase)

	var contactUseCase contact.ContactUseCase
	contactUseCase = contact.NewContactUseCase(config.DB, config.Log, config.Validate, contactRepository, contactProducer, slackClient)
	contactUseCase = contact.NewContactUseCaseMwLogger(config.Log, contactUseCase)

	var addressUseCase address.AddressUseCase
	addressUseCase = address.NewAddressUseCase(config.DB, config.Log, config.Validate, contactRepository, addressRepository, addressProducer, paymentClient)
	addressUseCase = address.NewAddressUseCaseMwLogger(config.Log, addressUseCase)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	contactController := http.NewContactController(contactUseCase, config.Log)
	addressController := http.NewAddressController(addressUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)
	traceIDMiddleware := middleware.NewTraceID()

	routeConfig := route.RouteConfig{
		App:               config.App,
		UserController:    userController,
		ContactController: contactController,
		AddressController: addressController,
		AuthMiddleware:    authMiddleware,
		TraceIDMiddleware: traceIDMiddleware,
	}

	routeConfig.Setup()
}
