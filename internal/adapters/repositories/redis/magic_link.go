package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hyperzoop/internal/core/entities"
	"time"

	"github.com/redis/go-redis/v9"
)

type MagicLinkRedisRepository struct {
	redis *redis.Client
}

func NewMagicLinkRedisRepository(redis *redis.Client) *MagicLinkRedisRepository {
	return &MagicLinkRedisRepository{
		redis: redis,
	}
}

func (r *MagicLinkRedisRepository) Create(link *entities.MagicLink) error {
	var bytes, err = json.Marshal(link)
	if err != nil {
		fmt.Println(err)
	}
	_, err = r.redis.Set(context.Background(), link.Code, string(bytes), time.Duration(time.Minute*15)).Result()
	return err
}

func (r *MagicLinkRedisRepository) FindValidByCode(code, cookie string) (*entities.MagicLink, error) {
	out, err := r.redis.Get(context.Background(), code).Result()
	if err != nil {
		return nil, err
	}
	if out == "" {
		return nil, errors.New("magic link not found")
	}
	var link entities.MagicLink
	if err := json.Unmarshal([]byte(out), &link); err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *MagicLinkRedisRepository) Invalidate(code string) error {
	_, err := r.redis.Del(context.Background(), code).Result()
	return err
}

func (r *MagicLinkRedisRepository) Update(link *entities.MagicLink) error {
	_, err := r.redis.Set(context.Background(), link.Code, link, 0).Result()
	return err
}
