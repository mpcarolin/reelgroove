package handlers

import (
	"context"
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/services"
	"github.com/mpcarolin/cinematch-server/internal/ui"
	"github.com/mpcarolin/cinematch-server/internal/ui/pages"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

func Summary(c echo.Context) error {
	// Parse query params
	movieId, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid movie id")
	}

	userLikes, err := utils.GetUserLikesFromCookie(c);
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting user likes from cookie")
	}

	recommendations, err := services.GetMovieRecommendations(movieId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting recommendations")
	}

	filteredRecommendations := []models.Movie{}
	for _, recommendation := range recommendations.Results {
		if slices.Contains(userLikes, strconv.Itoa(recommendation.Id)) {
			filteredRecommendations = append(filteredRecommendations, recommendation)
		}
	}

	component := pages.Summary(filteredRecommendations)
	page := ui.Page(component)
	return page.Render(context.Background(), c.Response().Writer)
}