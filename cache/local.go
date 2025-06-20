// Package cache implements a local LRU cache with singleflight.
package cache

import (
	"context"
	"errors"
	"time"

	"github.com/dgraph-io/ristretto"
	"golang.org/x/sync/singleflight"
)

type localCache struct {
	cache      *ristretto.Cache
	group      singleflight.Group
	defaultTTL time.Duration
}

func NewLocal(maxCost int64, defaultTTL time.Duration) (Cache, error) {
	c, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: maxCost * 10,
		MaxCost:     maxCost,
		BufferItems: 64,
	})
	if err != nil {
		return nil, err
	}
	return &localCache{
		cache:      c,
		defaultTTL: defaultTTL,
	}, nil
}

func (c *localCache) Get(ctx context.Context, key string, fetch func() (any, error), ttl time.Duration) (any, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if val, found := c.cache.Get(key); found {
		return val, nil
	}
	v, err, _ := c.group.Do(key, func() (any, error) {
		val, err := fetch()
		if err != nil {
			return nil, err
		}
		_ = c.cache.SetWithTTL(key, val, 1, ttl)
		return val, nil
	})
	return v, err
}

func (c *localCache) Set(ctx context.Context, key string, val any, ttl time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if key == "" {
		return errors.New("key cannot be empty")
	}
	c.cache.SetWithTTL(key, val, 1, ttl)
	return nil
}

func (c *localCache) Delete(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	c.cache.Del(key)
	return nil
}
