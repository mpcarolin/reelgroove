package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	store "github.com/eko/gocache/lib/v4/store"
	"github.com/mpcarolin/cinematch-server/internal/models"
)

// Movie-related structs have been moved to models/movie.go

func MovieMeetsUsageCriteria(movie models.Movie) bool {
	return movie.Poster != "" && movie.Popularity > 1 && movie.VoteAverage > 2 && movie.VoteCount > 25
}

func SearchMovies(search string) (models.MovieResponse, error) {
	cache := GetCache()
	cacheKey := fmt.Sprintf("search_movies_%s", search)
	if cachedResponse, cacheErr := cache.Get(context.Background(), cacheKey); cacheErr == nil {
		var response models.MovieResponse
		err := json.Unmarshal([]byte(cachedResponse), &response)
		if err != nil {
			slog.Error("error unmarshalling cached search response", "error", err)
			return models.MovieResponse{}, err
		}
		slog.Info("cache hit", "key", cacheKey, "value", response)
		return response, nil
	} else {
		slog.Error("error getting cached search response", "error", cacheErr)
	}

	var response models.MovieResponse
	err := json.Unmarshal([]byte(MockSearchResponse), &response)
	if err != nil {
		slog.Error("error unmarshalling mock search response", "error", err)
		return models.MovieResponse{}, err
	}

	// remove movies with no poster, low popularity, low vote average, low vote count, no video, or is adult
	// TODO: might look into filtering these out at the request level
	filteredResults := []models.Movie{}
	for _, movie := range response.Results {
		if MovieMeetsUsageCriteria(movie) { 
			filteredResults = append(filteredResults, movie)
		}
	}

	response.Results = filteredResults

	// TODO: could make this a go routine?
	serializedResponse, err := json.Marshal(response)
	if err != nil {
		slog.Error("error marshalling search response for storing in cache", "error", err)
		return models.MovieResponse{}, err
	}
	cacheSetErr := cache.Set(
		context.Background(),
		cacheKey,
		string(serializedResponse),
		store.WithExpiration(24*time.Hour),
	)
	if cacheSetErr != nil {
		slog.Error("error setting cache", "error", cacheSetErr)
	} else {
		slog.Info("cache set", "key", cacheKey, "value", string(serializedResponse))
	}

	return response, nil
}

func GetMovieRecommendations(movieId int) (models.RecommendationResponse, error) {
	recommendations := models.RecommendationResponse{}
	json.Unmarshal([]byte(MockRecommendationsResponse), &recommendations)

	filteredResults := []models.Movie{}
	for _, movie := range recommendations.Results {
		if MovieMeetsUsageCriteria(movie) { 
			filteredResults = append(filteredResults, movie)
		}
	}

	recommendations.Results = filteredResults[:10]

	return recommendations, nil
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

