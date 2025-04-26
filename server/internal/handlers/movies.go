package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/components"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

func cleanUpSearchQuery(search string) string {
	return strings.TrimSpace(strings.ToLower(search))
}

func GetMovieSearchResults(c echo.Context) error {
	search := c.QueryParam("search")
	searchQuery := cleanUpSearchQuery(search)
	ctx := c.(*models.RequestContext)

	slog.Info("GetMovieSearchResults", "searchQuery", searchQuery)

	response, err := utils.SearchMoviesCached(ctx.Cache, searchQuery)
	if err != nil {
		slog.Error("Error fetching movie search results", "error", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	movies := response.Results
	searchResults := []models.Movie{}
	for _, movie := range movies {
		if strings.Contains(strings.ToLower(movie.Title), searchQuery) {
			searchResults = append(searchResults, movie)
		}
	}

	component := components.MovieResults(searchResults)

	return component.Render(context.Background(), c.Response().Writer)
}
