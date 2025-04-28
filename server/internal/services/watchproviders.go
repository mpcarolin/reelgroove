package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

func GetWatchProvidersCached(cache *cache.Cache[string], movieId int) (*models.WatchProviders, error) {
	cacheKey := fmt.Sprintf("recommendations_watchproviders_%d", movieId)
	fetch := func() (*models.WatchProviders, error) {
		return GetWatchProviders(movieId)
	}
	return utils.WithCache(cache, cacheKey, fetch, 24*time.Hour);
}

func GetWatchProviders(movieId int) (*models.WatchProviders, error) {
	// response, err := requestGetWatchProviders(movieId)
	// if (err != nil) {
	// 	return nil, err
	// }
	response := MockWatchProvidersResponse

	var watchProvidersResponse models.WatchProvidersResponse
	err := json.Unmarshal([]byte(response), &watchProvidersResponse)
	if err != nil {
		slog.Error("error unmarshalling mock provider response", "error", err)
		return nil, err
	}

	// TODO: make this dynamic?
	country := "US" 
	watchProviders := watchProvidersResponse.Results[country]

	slices.SortFunc(watchProviders.Flatrate, func(a, b models.WatchProviderOption) int {
		return a.DisplayPriority - b.DisplayPriority
	})

	slog.Info("get watch providers for recommendation " + strconv.Itoa(movieId), "length=", watchProviders)

	return &watchProviders, nil
}

func requestGetWatchProviders(search string) (string, error) {
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
