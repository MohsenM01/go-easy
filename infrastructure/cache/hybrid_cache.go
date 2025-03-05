package cache

import (
	"go-easy/config"
	"time"
)

type HybridCache struct {
	primary   Cache
	secondary Cache
}

func NewHybridCache() *HybridCache {
	cfg := config.LoadConfig()

	inMemoryCache := newInMemoryCache()

	redisCache := newRedisCache(cfg.RedisConnection, "", 0)
	return &HybridCache{
		primary:   inMemoryCache,
		secondary: redisCache,
	}
}

func (h *HybridCache) Set(key string, value any, expiration time.Duration) error {
	if err := h.primary.Set(key, value, expiration); err != nil {
		return h.secondary.Set(key, value, expiration)
	}
	return nil
}

func (h *HybridCache) Get(key string) (any, error) {
	value, err := h.primary.Get(key)
	if err == nil {
		return value, nil
	}
	return h.secondary.Get(key)
}

func (h *HybridCache) Delete(key string) error {
	if err := h.primary.Delete(key); err != nil {
		return h.secondary.Delete(key)
	}
	return nil
}
