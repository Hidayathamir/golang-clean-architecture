package provider

import (
	"github.com/Hidayathamir/golang-clean-architecture/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:       cfg.GetRedisAddr(),
		Username:   cfg.GetRedisUsername(),
		Password:   cfg.GetRedisPassword(),
		DB:         cfg.GetRedisDB(),
		MaxRetries: 3,
	})
}
