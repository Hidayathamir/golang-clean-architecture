package repository

import "github.com/spf13/viper"

//go:generate moq -out=../mock/ImageRepository.go -pkg=mock . ImageRepository

type ImageRepository interface {
}

var _ ImageRepository = &ImageRepositoryImpl{}

type ImageRepositoryImpl struct {
	Config *viper.Viper
}

func NewImageRepository(cfg *viper.Viper) *ImageRepositoryImpl {
	return &ImageRepositoryImpl{
		Config: cfg,
	}
}
