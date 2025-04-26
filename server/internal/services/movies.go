package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

// TOOD: put somewhere else, it's a private variable but still accessible from all other files in this package,
// so the imports are suprising
var apiKey = os.Getenv("TMDB_API_KEY")

func SearchMoviesCached(cache *cache.Cache[string], search string) (*models.MovieSearchResponse, error) {
	cacheKey := fmt.Sprintf("movie_search_%s", utils.FormatCacheKey(search))
	fetch := func() (*models.MovieSearchResponse, error) {
		return SearchMovies(search)
	}
	return utils.WithCache(cache, cacheKey, fetch, 24*time.Hour);
}

func SearchMovies(search string) (*models.MovieSearchResponse, error) {
	response, err := requestSearchMovies(search)
	if (err != nil) {
		return nil, err
	}

	var searchResults models.MovieSearchResponse
	err = json.Unmarshal([]byte(response), &searchResults)
	if err != nil {
		slog.Error("error unmarshalling mock search response", "error", err)
		return nil, err
	}
	slog.Info("search movies response for query " + search, "length=", searchResults.TotalResults)

	return &searchResults, nil
}

func requestSearchMovies(search string) (string, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?query=%s&include_adult=false&language=en-US&page=1", url.QueryEscape(search))

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer " + apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("error sending search movies request", "error", err)
		return "", err
	}
	defer res.Body.Close()

	code := res.StatusCode;
	switch {
	case code >= 200 && code <= 299:
		body, err := io.ReadAll(res.Body)
		if err != nil {
			slog.Error("error reading search movies response", "error", err)
			return "", err
		}
		return string(body), nil
	case code == 400:
		msg := fmt.Sprintf("Bad input! Maybe query is bad. Query: %s. Code: %d, Message: %s", search, code, res.Status)
		slog.Error(msg)
		return "", errors.New("bad input")
	default:
		msg := fmt.Sprintf("Could not get movie search results from api. Code: %d, Message: %s", code, res.Status)
		slog.Error(msg)
		return "", errors.New("API did not provide movies for search query")
	}
}

// TODO: move to utils
func MovieMeetsUsageCriteria(movie models.Movie) bool {
	return movie.Poster != "" && movie.Popularity > 1 && movie.VoteAverage > 2 && movie.VoteCount > 25
}