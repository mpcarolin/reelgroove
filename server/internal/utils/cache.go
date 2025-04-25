package utils

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	lib_store "github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/redis/v4"
	redis "github.com/redis/go-redis/v9"
)

func GetCache() *cache.Cache[string] {
	redisHost := os.Getenv("REDIS_HOST")
	redisClient := redis.NewClient(&redis.Options{Addr: redisHost})
	redisStore := redis_store.NewRedis(redisClient, lib_store.WithExpiration(24*time.Hour))
	cacheManager := cache.New[string](redisStore)
	return cacheManager
}


func GetMarshalled[T any](cache *cache.Cache[string], cacheKey string) (*T, error) {
	res, err := cache.Get(context.Background(), cacheKey)
	if err != nil {
		return nil, err
	}

	var result T
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}