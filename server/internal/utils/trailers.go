package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/mpcarolin/cinematch-server/internal/models"
)

func GetBestMovieTrailerCached(cache *cache.Cache[string], movieId int) (*models.Trailer, error) {
	cacheKey := fmt.Sprintf("movie_trailer_%d", movieId)
	fetch := func(key string) (*models.Trailer, error) {
		return GetBestMovieTrailer(movieId)
	}
	return WithCache(cache, cacheKey, fetch, 24*time.Hour)
}

func GetBestMovieTrailer(movieId int) (*models.Trailer, error) {
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
		return &filteredTrailers[0], nil
	}

	return nil, errors.New("no trailers found")
}