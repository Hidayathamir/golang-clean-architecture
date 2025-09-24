package test

import (
	"context"
	"golang-clean-architecture/internal/config"
	"golang-clean-architecture/pkg/constant/configkey"
	"net"
	"os"
	"testing"

	"github.com/go-playground/validator/v10"
	gosqldrivermysql "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"gorm.io/gorm"
)

var app *fiber.App

var db *gorm.DB

var viperConfig *viper.Viper

var log *logrus.Logger

var validate *validator.Validate

// TestMain is the entry point for all tests in this package.
// It sets up global dependencies (logger, validator, Fiber app, DB, Kafka producer),
// starts a MySQL container, runs database migrations, bootstraps the application,
// executes test, and finally terminates the container before exiting.
func TestMain(m *testing.M) {
	viperConfig = config.NewViper()
	log = config.NewLogger(viperConfig)
	validate = config.NewValidator(viperConfig)
	app = config.NewFiber(viperConfig)
	mysqlContainer := newMysqlContainer(viperConfig)
	viperConfig.Set(configkey.DatabaseMigrations, "../db/migrations")
	config.Migrate(viperConfig, log)
	db = config.NewDatabase(viperConfig, log)
	producer := config.NewKafkaProducer(viperConfig, log)

	config.Bootstrap(&config.BootstrapConfig{
		DB:       db,
		App:      app,
		Log:      log,
		Validate: validate,
		Config:   viperConfig,
		Producer: producer,
	})

	code := m.Run()

	err := testcontainers.TerminateContainer(mysqlContainer)
	panicIfErr(err)

	os.Exit(code)
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func newMysqlContainer(viperConfig *viper.Viper) (mysqlContainer *mysql.MySQLContainer) {
	database := viperConfig.GetString(configkey.DatabaseName)
	username := viperConfig.GetString(configkey.DatabaseUsername)
	password := viperConfig.GetString(configkey.DatabasePassword)

	var err error
	mysqlContainer, err = mysql.Run(context.Background(),
		"mysql:8.0.36",
		mysql.WithDatabase(database),
		mysql.WithUsername(username),
		mysql.WithPassword(password),
	)
	panicIfErr(err)
	state, err := mysqlContainer.State(context.Background())
	panicIfErr(err)
	if !state.Running {
		panic("mysql container not running")
	}

	dbURL, err := mysqlContainer.ConnectionString(context.Background())
	panicIfErr(err)
	cfg, err := gosqldrivermysql.ParseDSN(dbURL)
	panicIfErr(err)

	_, port, err := net.SplitHostPort(cfg.Addr)
	panicIfErr(err)
	viperConfig.Set(configkey.DatabasePort, port)

	return mysqlContainer
}

func newKafkaContainer() (kafkaContainer *kafka.KafkaContainer) {
	var err error
	kafkaContainer, err = kafka.Run(context.Background(),
		"confluentinc/cp-kafka:7.2.15",
	)
	panicIfErr(err)
	state, err := kafkaContainer.State(context.Background())
	panicIfErr(err)
	if !state.Running {
		panic("kafka container not running")
	}
	return kafkaContainer
}
