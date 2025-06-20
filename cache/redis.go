// Package cache implements a Redis cache adapter.
package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisCache struct {
	rdb *redis.Client
}

func NewRedis(addr, password string, db int) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &redisCache{rdb: rdb}
}

func (c *redisCache) Get(ctx context.Context, key string, fetch func() (any, error), ttl time.Duration) (any, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if err == nil {
		return val, nil
	}
	if err != redis.Nil {
		return nil, err
	}
	// cache miss: fetch and set
	data, err := fetch()
	if err != nil {
		return nil, err
	}
	str, ok := data.(string)
	if !ok {
		return nil, errors.New("only string values are supported in redis cache")
	}
	_ = c.rdb.Set(ctx, key, str, ttl).Err()
	return str, nil
}

func (c *redisCache) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	str, ok := val.(string)
	if !ok {
		return errors.New("only string values are supported in redis cache")
	}
	return c.rdb.Set(ctx, key, str, ttl).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}
