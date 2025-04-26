package handlers

import (
	"context"
	"math"
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/components"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/services"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

// Renders the recommendations page for a single movie in the recommendations list
func GetSingleRecommendation(c echo.Context) error {
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

	// Render recommendations page
	return components.Page(
		components.Recommendations(models.TemplateContext{
			MovieId:         movieId,
			Trailer:         nextTrailer,
			Recommendations: recommendations.Results,
			UserLikes:       userLikes,
		}),
	).Render(context.Background(), c.Response().Writer)

}

// Fetches all the recommendations for a movie, and redirects to the first recommendation
func GetRecommendations(c echo.Context) error {
	ctx := c.(*models.RequestContext)

	movieId, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid movie id")
	}

	recommendations, err := services.GetMovieRecommendationsCached(ctx.Cache, movieId)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	recommendationIds := []int{}
	for _, movie := range recommendations.Results {
		recommendationIds = append(recommendationIds, movie.Id)
	}

	userLikesCookie := utils.CreateUserLikesCookie([]string{})
	http.SetCookie(c.Response().Writer, userLikesCookie)

	recommendationUrl := "/movie/" + strconv.Itoa(movieId) + "/recommendations/" + strconv.Itoa(recommendationIds[0])
	return c.Redirect(http.StatusSeeOther, recommendationUrl)
}

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

	nextRecommendationUrl := "/movie/" + strconv.Itoa(movieId) + "/recommendations/" + strconv.Itoa(nextRecommendationId)
	switch action {
	case "skip":
		userLikes = slices.DeleteFunc(userLikes, func(like string) bool { return like == strconv.Itoa(recommendationId) })
		userLikesCookie := utils.CreateUserLikesCookie(userLikes)
		http.SetCookie(c.Response().Writer, userLikesCookie)
		c.Response().Header().Set("HX-Redirect", nextRecommendationUrl)
		return c.NoContent(http.StatusOK)
	case "maybe":
		userLikes = append(userLikes, strconv.Itoa(recommendationId))
		userLikesCookie := utils.CreateUserLikesCookie(userLikes)
		http.SetCookie(c.Response().Writer, userLikesCookie)
		c.Response().Header().Set("HX-Redirect", nextRecommendationUrl)
		return c.NoContent(http.StatusOK)
	case "watch":
		// TODO: implement...
		return nil
	}

	return nil
}
