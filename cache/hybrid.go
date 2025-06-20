package cache

import (
	"context"
	"time"
)

type hybridCache struct {
	local Cache
	redis Cache
}

func NewHybridCache(local Cache, redis Cache) Cache {
	return &hybridCache{local: local, redis: redis}
}

func (c *hybridCache) Get(ctx context.Context, key string, fetch func() (any, error), ttl time.Duration) (any, error) {
	return c.local.Get(ctx, key, func() (any, error) {
		return c.redis.Get(ctx, key, fetch, ttl)
	}, ttl)
}

func (c *hybridCache) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	_ = c.local.Set(ctx, key, val, ttl)
	return c.redis.Set(ctx, key, val, ttl)
}

func (c *hybridCache) Delete(ctx context.Context, key string) error {
	_ = c.local.Delete(ctx, key)
	return c.redis.Delete(ctx, key)
}
