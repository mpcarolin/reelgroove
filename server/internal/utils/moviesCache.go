package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	store "github.com/eko/gocache/lib/v4/store"
	"github.com/mpcarolin/cinematch-server/internal/models"
)

func GetMovieSearchResultsCacheKey(search string) string {
	return fmt.Sprintf("search_movies_%s", search)
}

func GetMovieSearchFromCache(cache *cache.Cache[string], search string) (*models.MovieSearchResponse, error) {
	cacheKey := GetMovieSearchResultsCacheKey(search)
	cachedResponse, err := cache.Get(context.Background(), cacheKey);
	if cachedResponse != "" && err == nil {
		slog.Info("cache hit", "key", cacheKey)
		var response models.MovieSearchResponse
		err := json.Unmarshal([]byte(cachedResponse), &response)
		if err != nil {
			slog.Error("error unmarshalling cached search response", "error", err)
			return nil, err
		}
		return &response, nil
	}
	slog.Info("cache miss", "key", cacheKey)
	return nil, errors.New("no cached search response found")
}

func StoreMovieSearchResultsInCache(cache *cache.Cache[string], search string, response *models.MovieSearchResponse) error {
	cacheKey := GetMovieSearchResultsCacheKey(search)
	serializedResponse, err := json.Marshal(*response)
	if err != nil {
		slog.Error("error marshalling search response for storing in cache", "error", err)
		return err
	}
	cache.Set(context.Background(), cacheKey, string(serializedResponse), store.WithExpiration(24*time.Hour))
	slog.Info("cache set", "key", cacheKey)
	return nil
}