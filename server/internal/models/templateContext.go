package models

import (
	"math"
	"slices"
	"strconv"
)

type TemplateContext struct {
	MovieId         int
	Trailer         *Trailer
	Recommendations []Movie
	UserLikes       []string // id of recommendations the user has liked
	Autoplay        bool
}

// TODO: move these
func GetSkipUrl(context TemplateContext) string {
	return "/movie/" + strconv.Itoa(context.MovieId) + "/recommendations/" + strconv.Itoa(context.Trailer.MovieId) + "/skip"
}

func GetMaybeUrl(context TemplateContext) string {
	return "/movie/" + strconv.Itoa(context.MovieId) + "/recommendations/" + strconv.Itoa(context.Trailer.MovieId) + "/maybe"
}

func GetWatchUrl(context TemplateContext) string {
	return "/movie/" + strconv.Itoa(context.MovieId) + "/recommendations/" + strconv.Itoa(context.Trailer.MovieId) + "/watch"
}

func GetRecommendationUrl(movieId int, recommendationId int, autoplay *bool) string {
	queryString := ""
	if autoplay != nil && *autoplay {
		queryString = "?autoplay=on"
	}

	return "/movie/" + strconv.Itoa(movieId) + "/recommendations/" + strconv.Itoa(recommendationId) + queryString
}

// TODO: move all these to utils or something, and also reuse this one with the skip/maybe action?
func GetNextRecommendationUrl(context TemplateContext) string {
	currentRecommendationIndex := slices.IndexFunc(context.Recommendations, func(recommendation Movie) bool { return recommendation.Id == context.Trailer.MovieId })
	nextRecommendationIndex := math.Min(float64(currentRecommendationIndex+1), float64(len(context.Recommendations)-1))
	nextRecommendationId := context.Recommendations[int(nextRecommendationIndex)].Id

	// TODO: look up some go utility for query param string building
    autoplay := context.Autoplay
    queryString := ""
    if autoplay {
        queryString = "?autoplay=on"
    }
	return "/movie/" + strconv.Itoa(context.MovieId) + "/recommendations/" + strconv.Itoa(nextRecommendationId) + queryString
}
