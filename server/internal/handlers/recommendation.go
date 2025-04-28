package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/services"
	"github.com/mpcarolin/cinematch-server/internal/ui"
	"github.com/mpcarolin/cinematch-server/internal/ui/pages"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

// Renders the recommendations page for a single movie in the recommendations list
func GetRecommendationById(c echo.Context) error {
	ctx := c.(*models.RequestContext)

	// Parse query params
	movieId, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid movie id")
	}

	recommendationId, err := strconv.Atoi(c.Param("recommendationId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid recommendation id")
	}

	// Call api endpoints for movie data
	recommendations, err := services.GetMovieRecommendationsCached(ctx.Cache, movieId)
	if err != nil {
		// TODO: if there are no recommendations, we should display a ui indicating user should try a different base movie
		return c.String(http.StatusInternalServerError, err.Error())
	}

	nextTrailer, err := services.GetBestMovieTrailerCached(ctx.Cache, recommendationId)
	if err != nil {
		// TODO: in this case, there are no trailers for the movie, so we need to
		// perhaps remove this recommendation from the list, but not return this error here
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Get user likes from cookie, to ensure page is rendered with correct user likes
	userLikes, err := utils.GetUserLikesFromCookie(c)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	autoplay := false
	if ctx.QueryParam("autoplay") == "on" {
		autoplay = true
	}

	recommendationViewModel := pages.RecommendationViewModel{
		MovieId:         movieId,
		CurrentRecommendationId: recommendationId,
		Recommendations: recommendations.Results,
		UserLikes:       userLikes,
		Trailer:         nextTrailer,
		Settings:        models.RecommendationSettings{Autoplay: autoplay},
	}

	return ui.Page(
		pages.Recommendation(&recommendationViewModel),
	).Render(context.Background(), c.Response().Writer)

}
