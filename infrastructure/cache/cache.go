package cache

import "time"

type Cache interface {
    Set(key string, value any, expiration time.Duration) error
    Get(key string) (any, error)
    Delete(key string) error
}
