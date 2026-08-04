package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/constant/cachekey"
	"github.com/Hidayathamir/golang-clean-architecture/pkg/errkit"
	"github.com/redis/go-redis/v9"
)

//go:generate moq -out=../../mock/MockCacheUser.go -pkg=mock . UserCache

type UserCache interface {
	Get(ctx context.Context, id int64) (*entity.User, error)
	Set(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int64) error
}

type UserCacheImpl struct {
	client *redis.Client
}

var _ UserCache = &UserCacheImpl{}

func NewUserCache(client *redis.Client) UserCache {
	return &UserCacheImpl{
		client: client,
	}
}

func (c *UserCacheImpl) Get(ctx context.Context, id int64) (*entity.User, error) {
	val, err := c.client.Get(ctx, cachekey.UserID(id)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, errkit.AddFuncName(err, "cache.(*UserCacheImpl).Get")
	}

	user := entity.User{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, errkit.AddFuncName(err, "cache.(*UserCacheImpl).Get")
	}

	return &user, nil
}

func (c *UserCacheImpl) Set(ctx context.Context, user *entity.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return errkit.AddFuncName(err, "cache.(*UserCacheImpl).Set")
	}

	err = c.client.Set(ctx, cachekey.UserID(user.ID), data, 1*time.Hour).Err()
	if err != nil {
		return errkit.AddFuncName(err, "cache.(*UserCacheImpl).Set")
	}

	return nil
}

func (c *UserCacheImpl) Delete(ctx context.Context, id int64) error {
	err := c.client.Del(ctx, cachekey.UserID(id)).Err()
	if err != nil {
		return errkit.AddFuncName(err, "cache.(*UserCacheImpl).Delete")
	}
	return nil
}
