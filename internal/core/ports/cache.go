package ports

import "time"

type RedisCacheRepository interface {
	Set(key string, value string, expiration time.Duration) error
	Get(key string) (string, error)
	Invalidate(key string) error
}
