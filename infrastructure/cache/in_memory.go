package cache

import (
	"errors"
	"sync"
	"time"
)

type inMemoryCache struct {
	data map[string]cacheItem
	mu   sync.RWMutex
}

type cacheItem struct {
	value      any
	expiration time.Time
}

func newInMemoryCache() Cache {
	return &inMemoryCache{
		data: make(map[string]cacheItem),
	}
}

func (c *inMemoryCache) Set(key string, value any, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(expiration),
	}
	return nil
}

func (c *inMemoryCache) Get(key string) (any, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, exists := c.data[key]
	if !exists || time.Now().After(item.expiration) {
		return nil, errors.New("key not found or expired")
	}
	return item.value, nil
}

func (c *inMemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}
