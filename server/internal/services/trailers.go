package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

func GetBestMovieTrailerCached(cache *cache.Cache[string], movieId int) (*models.Trailer, error) {
	cacheKey := fmt.Sprintf("movie_trailer_%d", movieId)
	fetch := func() (*models.Trailer, error) {
		return GetBestMovieTrailer(movieId)
	}
	return utils.WithCache(cache, cacheKey, fetch, 24*time.Hour)
}

func GetBestMovieTrailer(movieId int) (*models.Trailer, error) {
	response, err := requestTrailers(movieId)
	if (err != nil) {
		return nil, err
	}
	trailerResponse := models.TrailerResponse{}
	json.Unmarshal([]byte(response), &trailerResponse)

	filteredTrailers := []models.Trailer{}
	for _, trailer := range trailerResponse.Results {
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

func requestTrailers(movieId int) (string, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/videos?language=en-US", movieId)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer " + apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("error sending trailers request", "error", err)
		return "", err
	}
	defer res.Body.Close()

	code := res.StatusCode;
	switch {
	case code >= 200 && code <= 299:
		body, err := io.ReadAll(res.Body)
		if err != nil {
			slog.Error("error reading trailers response", "error", err)
			return "", err
		}
		return string(body), nil
	case code == 400:
		msg := fmt.Sprintf("Bad input! Maybe movieId is bad. movieId: %d. Code: %d, Message: %s", movieId, code, res.Status)
		slog.Error(msg)
		return "", errors.New("bad input")
	default:
		msg := fmt.Sprintf("Could not get movie trailers from api. Code: %d, Message: %s", code, res.Status)
		slog.Error(msg)
		return "", errors.New("API did not provide trailers for movie")
	}
}
