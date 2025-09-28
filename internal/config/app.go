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
	userRepository = repository.NewUserRepositoryMwLogger(userRepository)

	var contactRepository repository.ContactRepository
	contactRepository = repository.NewContactRepository(config.Log)
	contactRepository = repository.NewContactRepositoryMwLogger(contactRepository)

	var addressRepository repository.AddressRepository
	addressRepository = repository.NewAddressRepository(config.Log)
	addressRepository = repository.NewAddressRepositoryMwLogger(addressRepository)

	// setup producer
	var userProducer messaging.UserProducer
	userProducer = messaging.NewUserProducer(config.Producer, config.Log)
	userProducer = messaging.NewUserProducerMwLogger(userProducer)

	var contactProducer messaging.ContactProducer
	contactProducer = messaging.NewContactProducer(config.Producer, config.Log)
	contactProducer = messaging.NewContactProducerMwLogger(contactProducer)

	var addressProducer messaging.AddressProducer
	addressProducer = messaging.NewAddressProducer(config.Producer, config.Log)
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
	userUsecase = user.NewUserUsecase(config.DB, config.Log, config.Validate, userRepository, userProducer, s3Client, slackClient)
	userUsecase = user.NewUserUsecaseMwLogger(userUsecase)

	var contactUsecase contact.ContactUsecase
	contactUsecase = contact.NewContactUsecase(config.DB, config.Log, config.Validate, contactRepository, contactProducer, slackClient)
	contactUsecase = contact.NewContactUsecaseMwLogger(contactUsecase)

	var addressUsecase address.AddressUsecase
	addressUsecase = address.NewAddressUsecase(config.DB, config.Log, config.Validate, contactRepository, addressRepository, addressProducer, paymentClient)
	addressUsecase = address.NewAddressUsecaseMwLogger(addressUsecase)

	// setup controller
	userController := http.NewUserController(userUsecase, config.Log)
	contactController := http.NewContactController(contactUsecase, config.Log)
	addressController := http.NewAddressController(addressUsecase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUsecase)
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
