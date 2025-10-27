package user2

import (
	"context"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/model"
	"github.com/Hidayathamir/golang-clean-architecture/internal/repository"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/configkey"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

//go:generate moq -out=../../mock/User2Usecase.go -pkg=mock . User2Usecase

type User2Usecase interface {
	Register(ctx context.Context, req *model.RegisterUser2Request) (*model.User2Response, error)
	Login(ctx context.Context, req *model.LoginUser2Request) (*model.User2TokenResponse, error)
	Profile(ctx context.Context, req *model.GetUser2Request) (*model.User2Response, error)
	VerifyToken(ctx context.Context, req *model.VerifyUser2TokenRequest) (*model.User2Auth, error)
}

var _ User2Usecase = &User2UsecaseImpl{}

type User2UsecaseImpl struct {
	Config   *viper.Viper
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate

	User2Repository repository.User2Repository

	jwtSecret []byte
	jwtTTL    time.Duration
}

func NewUser2Usecase(
	cfg *viper.Viper,
	log *logrus.Logger,
	db *gorm.DB,
	validate *validator.Validate,
	user2Repository repository.User2Repository,
) *User2UsecaseImpl {
	secret := cfg.GetString(configkey.AuthJWTSecret)
	if secret == "" {
		log.Warn("auth.jwt.secret is empty, falling back to development default")
		secret = "change-me"
	}

	ttlSeconds := cfg.GetInt(configkey.AuthJWTTTLSeconds)
	if ttlSeconds <= 0 {
		ttlSeconds = 3600
	}

	return &User2UsecaseImpl{
		Config:          cfg,
		Log:             log,
		DB:              db,
		Validate:        validate,
		User2Repository: user2Repository,
		jwtSecret:       []byte(secret),
		jwtTTL:          time.Duration(ttlSeconds) * time.Second,
	}
}
