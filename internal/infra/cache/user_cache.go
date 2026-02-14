package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Hidayathamir/golang-clean-architecture/internal/entity"
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

func (c *UserCacheImpl) getKey(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func (c *UserCacheImpl) Get(ctx context.Context, id int64) (*entity.User, error) {
	val, err := c.client.Get(ctx, c.getKey(id)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	user := entity.User{}
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, errkit.AddFuncName(err)
	}

	return &user, nil
}

func (c *UserCacheImpl) Set(ctx context.Context, user *entity.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return errkit.AddFuncName(err)
	}

	err = c.client.Set(ctx, c.getKey(user.ID), data, 1*time.Hour).Err()
	if err != nil {
		return errkit.AddFuncName(err)
	}

	return nil
}

func (c *UserCacheImpl) Delete(ctx context.Context, id int64) error {
	err := c.client.Del(ctx, c.getKey(id)).Err()
	if err != nil {
		return errkit.AddFuncName(err)
	}
	return nil
}
