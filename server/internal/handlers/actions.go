package handlers

import (
	"context"
	"math"
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/services"
	"github.com/mpcarolin/cinematch-server/internal/ui/pages"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

func LikeRecommendation(c echo.Context) error {
	recommendationViewModel, err := initRecommendationViewModel(c)
	if err != nil {
		return err
	}

	recommendationViewModel.UserLikes = append(recommendationViewModel.UserLikes, strconv.Itoa(recommendationViewModel.CurrentRecommendationId))
	userLikesCookie := utils.CreateUserLikesCookie(recommendationViewModel.UserLikes)
	http.SetCookie(c.Response().Writer, userLikesCookie)

	return pages.Recommendation(recommendationViewModel).Render(context.Background(), c.Response().Writer)
}

func SkipRecommendation(c echo.Context) error {
	recommendationViewModel, err := initRecommendationViewModel(c)
	if err != nil {
		return err
	}

	recommendationViewModel.UserLikes = slices.DeleteFunc(recommendationViewModel.UserLikes, func(like string) bool { return like == strconv.Itoa(recommendationViewModel.CurrentRecommendationId) })
	userLikesCookie := utils.CreateUserLikesCookie(recommendationViewModel.UserLikes)
	http.SetCookie(c.Response().Writer, userLikesCookie)

	return pages.Recommendation(recommendationViewModel).Render(context.Background(), c.Response().Writer)
}


// Initializes a recommendation view model, given the path and query params of the current request url,
// and the user's likes from the cookie.
func initRecommendationViewModel(c echo.Context) (*pages.RecommendationViewModel, error) {
	ctx := c.(*models.RequestContext)

	movieId, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		return nil, c.String(http.StatusBadRequest, "Invalid movie id")
	}

	recommendationId, err := strconv.Atoi(c.Param("recommendationId"))
	if err != nil {
		return nil, c.String(http.StatusBadRequest, "Invalid recommendation id")
	}

	userLikes, err := utils.GetUserLikesFromCookie(c)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, err.Error())
	}

	recommendations, err := services.GetMovieRecommendationsCached(ctx.Cache, movieId)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, err.Error())
	}

	currentRecommendationIndex := slices.IndexFunc(recommendations.Results, func(recommendation models.Movie) bool { return recommendation.Id == recommendationId })
	nextRecommendationIndex := math.Min(float64(currentRecommendationIndex+1), float64(len(recommendations.Results)-1))
	nextRecommendationId := recommendations.Results[int(nextRecommendationIndex)].Id

	nextTrailer, err := services.GetBestMovieTrailerCached(ctx.Cache, nextRecommendationId)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, err.Error())
	}

	settings := models.RecommendationSettings{
		Autoplay: c.QueryParam("autoplay") == "on",
	}

	recommendationViewModel := pages.RecommendationViewModel{
		MovieId:                 movieId,
		CurrentRecommendationId: recommendationId,
		Recommendations:         recommendations.Results,
		UserLikes:               userLikes,
		Trailer:                 nextTrailer,
		Settings:                settings,
	}

	return &recommendationViewModel, nil
}
