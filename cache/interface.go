package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, fetch func() (any, error), ttl time.Duration) (any, error)
	Set(ctx context.Context, key string, val any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
