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

func GetRecommendationsCacheKey(movieId int) string {
	return fmt.Sprintf("movie_recommendations_by_movie_id_%d", movieId)
}

func GetRecommendationsFromCache(cache *cache.Cache[string], movieId int) (*models.RecommendationResponse, error) {
	cacheKey := GetRecommendationsCacheKey(movieId)
	cachedResponse, err := cache.Get(context.Background(), cacheKey);
	if cachedResponse != "" && err == nil {
		var response models.RecommendationResponse
		err := json.Unmarshal([]byte(cachedResponse), &response)
		if err != nil {
			slog.Error("error unmarshalling cached search response", "error", err)
			return nil, err
		} else {
			slog.Info("cache hit", "key", cacheKey)
		}
		return &response, nil
	}

	slog.Info("cache miss", "key", cacheKey)
	return nil, errors.New("no cached search response found")
}

func StoreRecommendationsInCache(cache *cache.Cache[string], movieId int, response *models.RecommendationResponse) error {
	cacheKey := GetRecommendationsCacheKey(movieId)
	serializedResponse, err := json.Marshal(*response)
	if err != nil {
		slog.Error("error marshalling search response for storing in cache", "error", err)
		return err
	}
	cache.Set(context.Background(), cacheKey, string(serializedResponse), store.WithExpiration(24*time.Hour))
	slog.Info("cache set", "key", cacheKey)
	return nil
}