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

// Handles the action of a recommendation, either skipping, maybe, or watching
// Updates the user likes cookie, and redirects to the next recommendation
func HandleRecommendationAction(c echo.Context) error {
	ctx := c.(*models.RequestContext)

	movieId, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid movie id")
	}

	recommendationId, err := strconv.Atoi(c.Param("recommendationId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid recommendation id")
	}

	action := c.Param("action")

	userLikes, err := utils.GetUserLikesFromCookie(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	recommendations, err := services.GetMovieRecommendationsCached(ctx.Cache, movieId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	currentRecommendationIndex := slices.IndexFunc(recommendations.Results, func(recommendation models.Movie) bool { return recommendation.Id == recommendationId })
	nextRecommendationIndex := math.Min(float64(currentRecommendationIndex+1), float64(len(recommendations.Results)-1))
	nextRecommendationId := recommendations.Results[int(nextRecommendationIndex)].Id

	switch action {
	case "skip":
		userLikes = slices.DeleteFunc(userLikes, func(like string) bool { return like == strconv.Itoa(recommendationId) })
		userLikesCookie := utils.CreateUserLikesCookie(userLikes)
		http.SetCookie(c.Response().Writer, userLikesCookie)
		nextTrailer, err := services.GetBestMovieTrailerCached(ctx.Cache, nextRecommendationId)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return pages.Recommendations(models.TemplateContext{
			MovieId:         movieId,
			Trailer:         nextTrailer,
			Recommendations: recommendations.Results,
			UserLikes:       userLikes,
		}).Render(context.Background(), c.Response().Writer)
	case "maybe":
		userLikes = append(userLikes, strconv.Itoa(recommendationId))
		userLikesCookie := utils.CreateUserLikesCookie(userLikes)
		http.SetCookie(c.Response().Writer, userLikesCookie)
		nextTrailer, err := services.GetBestMovieTrailerCached(ctx.Cache, nextRecommendationId)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return pages.Recommendations(models.TemplateContext{
			MovieId:         movieId,
			Trailer:         nextTrailer,
			Recommendations: recommendations.Results,
			UserLikes:       userLikes,
		}).Render(context.Background(), c.Response().Writer)
	case "watch":
		// TODO: implement...
		return nil
	}

	return nil
}
