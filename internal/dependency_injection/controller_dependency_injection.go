package dependency_injection

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http"
)

type Controllers struct {
	UserController  *http.UserController
	ImageController *http.ImageController
}

func SetupControllers(cfg *config.Config, usecases *Usecases) *Controllers {
	userController := http.NewUserController(cfg, usecases.UserUsecase)
	imageController := http.NewImageController(cfg, usecases.ImageUsecase)

	return &Controllers{
		UserController:  userController,
		ImageController: imageController,
	}
}
