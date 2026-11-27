package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"GrandExchange/internal/models"
	"GrandExchange/internal/persist"

	"github.com/redis/go-redis/v9"
)

const cacheTTL = 30 * time.Second

type CacheService interface {
	Get(context.Context, string) (*models.SearchResponse, error)
	Set(context.Context, string, models.SearchResponse) error
}

type OsrsCache struct {
	redisCtx persist.RedisContext
}

func NewCacheService(redisCtx persist.RedisContext) CacheService {
	return &OsrsCache{redisCtx}
}

func (c *OsrsCache) Get(ctx context.Context, key string) (*models.SearchResponse, error) {
	value, err := c.redisCtx.Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("could not get from cache: %w", err)
	}

	var result models.SearchResponse
	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal cache value: %w", err)
	}
	return &result, nil
}

func (c *OsrsCache) Set(ctx context.Context, key string, value models.SearchResponse) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("could not marshal data for cache: %w", err)
	}

	if err := c.redisCtx.Rdb.Set(ctx, key, data, cacheTTL).Err(); err != nil {
		return fmt.Errorf("could not set cache: %w", err)
	}

	return nil
}
