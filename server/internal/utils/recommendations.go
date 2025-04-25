package utils

import (
	"encoding/json"
	"errors"
	"slices"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/mpcarolin/cinematch-server/internal/models"
)

func GetMovieRecommendationsCached(cache *cache.Cache[string], movieId int) (*models.RecommendationResponse, error) {
	cachedResponse, err := GetRecommendationsFromCache(cache, movieId)
	if cachedResponse != nil && err == nil {
		return cachedResponse, nil
	}

	response, err := GetMovieRecommendations(movieId)
	if err != nil {
		return nil, err
	}

	go StoreRecommendationsInCache(cache, movieId, response)

	return response, nil
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

func GetBestMovieTrailer(movieId int) (models.Trailer, error) {
	trailers := []models.Trailer{}
	json.Unmarshal([]byte(MockTrailersResponse), &trailers)

	filteredTrailers := []models.Trailer{}
	for _, trailer := range trailers {
		if trailer.Site == "YouTube" {
			trailer.MovieId = movieId
			filteredTrailers = append(filteredTrailers, trailer)
		}
	}

	slices.SortFunc(filteredTrailers, func(a, b models.Trailer) int {
		if a.Type == "Trailer" && b.Type != "Trailer" {
			return -1;
		} else if a.Type != "Trailer" && b.Type == "Trailer" {
			return 1;
		} else if a.Official && !b.Official {
			return -1;
		} else if !a.Official && b.Official {
			return 1;
		} else {
			return 0;
		}
	})

	if len(filteredTrailers) > 0 {
		return filteredTrailers[0], nil
	}

	return models.Trailer{}, errors.New("no trailers found")
}
