package persist

import (
	"GrandExchange/internal/config"

	"github.com/redis/go-redis/v9"
)

type RedisContext struct {
	Rdb *redis.Client
}

func NewRedisContext() RedisContext {
	return RedisContext{redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: "",
		DB:       0,
	})}
}
