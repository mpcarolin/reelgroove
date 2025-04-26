package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/mpcarolin/cinematch-server/internal/models"
)

func MovieMeetsUsageCriteria(movie models.Movie) bool {
	return movie.Poster != "" && movie.Popularity > 1 && movie.VoteAverage > 2 && movie.VoteCount > 25
}

func SearchMoviesCached(cache *cache.Cache[string], search string) (*models.MovieSearchResponse, error) {
	cacheKey := fmt.Sprintf("movie_search_%s", search)
	fetch := func(searchKey string) (*models.MovieSearchResponse, error) {
		return SearchMovies(searchKey)
	}
	return WithCache(cache, cacheKey, fetch, 24*time.Hour);
}

func SearchMovies(search string) (*models.MovieSearchResponse, error) {
	var response models.MovieSearchResponse
	err := json.Unmarshal([]byte(MockSearchResponse), &response)
	if err != nil {
		slog.Error("error unmarshalling mock search response", "error", err)
		return nil, err
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

	time.Sleep(2000 * time.Millisecond) // 2 second delay for testing

	return &response, nil
}
