package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/services"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

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

	recommendationUrl := "/movie/" + strconv.Itoa(movieId) + "/recommendations/" + strconv.Itoa(recommendationIds[0]) + "?autoplay=on"
	return c.Redirect(http.StatusSeeOther, recommendationUrl)
}

