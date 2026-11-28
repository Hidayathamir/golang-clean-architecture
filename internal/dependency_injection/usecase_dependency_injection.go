package dependency_injection

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/chatusecase"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/idempotencyusecase"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/imageusecase"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/notifusecase"
	"github.com/Hidayathamir/golang-clean-architecture/internal/usecase/userusecase"
	"gorm.io/gorm"
)

type Usecases struct {
	UserUsecase        userusecase.UserUsecase
	ImageUsecase       imageusecase.ImageUsecase
	NotifUsecase       notifusecase.NotifUsecase
	IdempotencyUsecase idempotencyusecase.IdempotencyUsecase
	ChatUsecase        chatusecase.ChatUsecase
}

func SetupUsecases(
	cfg *config.Config,
	db *gorm.DB,
	caches *Caches,
	messagings *Messagings,
	repositories *Repositories,
	searches *Searches,
	storages *Storages,
	emails *Emails,
) *Usecases {
	var idempotencyUsecase idempotencyusecase.IdempotencyUsecase
	idempotencyUsecase = idempotencyusecase.NewIdempotencyUsecase(cfg, db, repositories.IdempotencyRepository)
	idempotencyUsecase = idempotencyusecase.NewIdempotencyUsecaseMwLogger(idempotencyUsecase)

	var notifUsecase notifusecase.NotifUsecase
	notifUsecase = notifusecase.NewNotifUsecase(cfg, db, repositories.NotificationRepository, repositories.UserRepository, caches.NotificationPublisher)
	notifUsecase = notifusecase.NewNotifUsecaseMwLogger(notifUsecase)

	var userUsecase userusecase.UserUsecase
	userUsecase = userusecase.NewUserUsecase(cfg, db, repositories.UserRepository, repositories.UserStatRepository, repositories.FollowRepository, messagings.UserProducer, messagings.NotifProducer, messagings.EmailProducer, storages.S3Client, emails.Sender, caches.UserCache, notifUsecase)
	userUsecase = userusecase.NewUserUsecaseMwLogger(userUsecase)

	var imageUsecase imageusecase.ImageUsecase
	imageUsecase = imageusecase.NewImageUsecase(cfg, db, repositories.ImageRepository, repositories.LikeRepository, repositories.CommentRepository, repositories.FollowRepository, repositories.UserRepository, messagings.ImageProducer, messagings.NotifProducer, storages.S3Client, searches.ImageSearch, notifUsecase)
	imageUsecase = imageusecase.NewImageUsecaseMwLogger(imageUsecase)

	var chatUsecase chatusecase.ChatUsecase
	chatUsecase = chatusecase.NewChatUsecase(cfg, db, repositories.ChatRepository, repositories.UserRepository, caches.NotificationPublisher)
	chatUsecase = chatusecase.NewChatUsecaseMwLogger(chatUsecase)

	return &Usecases{
		UserUsecase:        userUsecase,
		ImageUsecase:       imageUsecase,
		NotifUsecase:       notifUsecase,
		IdempotencyUsecase: idempotencyUsecase,
		ChatUsecase:        chatUsecase,
	}
}
