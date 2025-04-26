package services

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

func GetMovieRecommendationsCached(cache *cache.Cache[string], movieId int) (*models.RecommendationResponse, error) {
	cacheKey := fmt.Sprintf("movie_recommendations_%d", movieId)
	fetch := func(key string) (*models.RecommendationResponse, error) {
		return GetMovieRecommendations(movieId)
	}
	return utils.WithCache(cache, cacheKey, fetch, 24*time.Hour);
}

func GetMovieRecommendations(movieId int) (*models.RecommendationResponse, error) {
	recommendations := models.RecommendationResponse{}
	json.Unmarshal([]byte(MockRecommendationsResponse), &recommendations)

	filteredResults := []models.Movie{}
	for _, movie := range recommendations.Results {
		if MovieMeetsUsageCriteria(movie) { 
			filteredResults = append(filteredResults, movie)
		}
	}

	recommendations.Results = filteredResults[:10]

	return &recommendations, nil
}
