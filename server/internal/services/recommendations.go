package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

func GetMovieRecommendationsCached(cache *cache.Cache[string], movieId int) (*models.RecommendationResponse, error) {
	cacheKey := fmt.Sprintf("movie_recommendations_%d", movieId)
	fetch := func() (*models.RecommendationResponse, error) {
		return GetMovieRecommendations(movieId)
	}
	return utils.WithCache(cache, cacheKey, fetch, 24*time.Hour);
}

func GetMovieRecommendations(movieId int) (*models.RecommendationResponse, error) {
	response, err := requestRecommendations(movieId)
	if (err != nil) {
		return nil, err
	}

	recommendations := models.RecommendationResponse{}
	json.Unmarshal([]byte(response), &recommendations)

	filteredResults := []models.Movie{}
	for _, movie := range recommendations.Results {
		if MovieMeetsUsageCriteria(movie) { 
			filteredResults = append(filteredResults, movie)
		}
	}

	recommendations.Results = filteredResults

	return &recommendations, nil
}

func requestRecommendations(movieId int) (string, error) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/recommendations?language=en-US&page=1", movieId)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer " + apiKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("error sending recommendations request", "error", err)
		return "", err
	}
	defer res.Body.Close()

	code := res.StatusCode;
	switch {
	case code >= 200 && code <= 299:
		body, err := io.ReadAll(res.Body)
		if err != nil {
			slog.Error("error reading recommendations response", "error", err)
			return "", err
		}
		return string(body), nil
	case code == 400:
		msg := fmt.Sprintf("Bad input! Maybe movieId is bad. movieId: %d. Code: %d, Message: %s", movieId, code, res.Status)
		slog.Error(msg)
		return "", errors.New("bad input")
	default:
		msg := fmt.Sprintf("Could not get movie recommendations from api. Code: %d, Message: %s", code, res.Status)
		slog.Error(msg)
		return "", errors.New("API did not provide recommendations for movie")
	}
}
