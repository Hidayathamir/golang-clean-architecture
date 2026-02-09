package dependency_injection

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/delivery/http"
	"github.com/spf13/viper"
)

type Controllers struct {
	UserController  *http.UserController
	ImageController *http.ImageController
}

func SetupControllers(viperConfig *viper.Viper, usecases *Usecases) *Controllers {
	userController := http.NewUserController(viperConfig, usecases.UserUsecase)
	imageController := http.NewImageController(viperConfig, usecases.ImageUsecase)

	return &Controllers{
		UserController:  userController,
		ImageController: imageController,
	}
}
