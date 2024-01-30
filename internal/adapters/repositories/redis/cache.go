package repositories

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCacheRepository struct {
	redis *redis.Client
}

func NewRedisCacheRepository(redis *redis.Client) *RedisCacheRepository {
	return &RedisCacheRepository{
		redis: redis,
	}
}

func (r *RedisCacheRepository) Set(key string, value string, expiration time.Duration) error {
	return r.redis.Set(context.Background(), key, value, expiration).Err()
}

func (r *RedisCacheRepository) Get(key string) (string, error) {
	return r.redis.Get(context.Background(), key).Result()
}

func (r *RedisCacheRepository) Invalidate(key string) error {
	return r.redis.Del(context.Background(), key).Err()
}
