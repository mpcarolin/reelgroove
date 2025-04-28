package models

import (
	"math"
	"slices"
	"strconv"
)

// TODO: move these to utils/urls.go

func GetUpdateSettingsUrl(movieId int, recommendationId int) string {
	return "/movie/" + strconv.Itoa(movieId) + "/recommendations/" + strconv.Itoa(recommendationId) + "/settings"
}

// TODO: move these
func GetSkipUrl(movieId int, recommendationId int) string {
	return "/movie/" + strconv.Itoa(movieId) + "/recommendations/" + strconv.Itoa(recommendationId) + "/skip"
}

func GetMaybeUrl(movieId int, recommendationId int) string {
	return "/movie/" + strconv.Itoa(movieId) + "/recommendations/" + strconv.Itoa(recommendationId) + "/maybe"
}

func GetWatchUrl(movieId int, recommendationId int) string {
	return "/movie/" + strconv.Itoa(movieId) + "/recommendations/" + strconv.Itoa(recommendationId) + "/watch"
}

func GetRecommendationUrl(movieId int, recommendationId int, autoplay *bool) string {
	queryString := ""
	if autoplay != nil && *autoplay {
		queryString = "?autoplay=on"
	}

	return "/movie/" + strconv.Itoa(movieId) + "/recommendations/" + strconv.Itoa(recommendationId) + queryString
}

// TODO: move all these to utils or something, and also reuse this one with the skip/maybe action?
func GetNextRecommendationUrl(movieId int, recommendations []Movie, currentRecommendationId int, autoplay *bool) string {
	currentRecommendationIndex := slices.IndexFunc(recommendations, func(recommendation Movie) bool { return recommendation.Id == currentRecommendationId })
	nextRecommendationIndex := math.Min(float64(currentRecommendationIndex+1), float64(len(recommendations)-1))
	nextRecommendationId := recommendations[int(nextRecommendationIndex)].Id

	// TODO: look up some go utility for query param string building
    queryString := ""
    if autoplay != nil && *autoplay {
        queryString = "?autoplay=on"
    }
	return "/movie/" + strconv.Itoa(movieId) + "/recommendations/" + strconv.Itoa(nextRecommendationId) + queryString
}
