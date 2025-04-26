package utils

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	lib_store "github.com/eko/gocache/lib/v4/store"
	redis_store "github.com/eko/gocache/store/redis/v4"
	redis "github.com/redis/go-redis/v9"
)

// CacheableData represents any function that fetches data and might return an error
type CacheableData[K comparable, V any] func() (V, error)

func Serialize[V any](data V) (string, error) {
	serializedData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(serializedData), nil
}

func Deserialize[V any](data string) (*V, error) {
	var result V
	err := json.Unmarshal([]byte(data), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func StoreInCache[K comparable, V any](
	cache *cache.Cache[string],
	key K,
	data *V,
	ttl time.Duration,
) error {
	serializedData, err := Serialize(*data)
	if err != nil {
		return err
	}
	return cache.Set(context.Background(), key, serializedData, lib_store.WithExpiration(ttl))
}

func GetFromCache[K comparable, V any](	
	cache *cache.Cache[string],
	key K,
) (*V, error) {
	cachedResponse, err := cache.Get(context.Background(), key);
	if cachedResponse != "" && err == nil {
		deserializedData, deserializeErr := Deserialize[V](cachedResponse)
		if deserializeErr != nil {
			slog.Error("error unmarshalling cached search response", "error", deserializeErr)
			return nil, deserializeErr
		}
		return deserializedData, nil
	}
	return nil, nil
}

// WithCache wraps any data fetching function with caching logic
func WithCache[K comparable, V any](
	cache *cache.Cache[string],
	key K,
	fetch CacheableData[K, *V],
	ttl time.Duration,
) (*V, error) {
	cachedData, err := GetFromCache[K, V](cache, key)
	if cachedData != nil && err == nil {
		return cachedData, nil
	} 

	// Cache miss - fetch fresh data
	slog.Info("cache miss", "key", key)
	data, err := fetch()
	if err != nil {
		slog.Error("error fetching data", "error", err)
		return nil, err
	}

	// Store in cache asynchronously
	go StoreInCache(cache, key, data, ttl)

	return data, nil
}

func FormatCacheKey(key string) string {
	return strings.ToLower(strings.ReplaceAll(strings.Trim(key, " "), " ", "_"))
}

func GetCache() *cache.Cache[string] {
	redisHost := os.Getenv("REDIS_HOST")
	redisClient := redis.NewClient(&redis.Options{Addr: redisHost})
	redisStore := redis_store.NewRedis(redisClient, lib_store.WithExpiration(24*time.Hour))
	cacheManager := cache.New[string](redisStore)
	return cacheManager
}

