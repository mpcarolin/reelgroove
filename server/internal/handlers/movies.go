package handlers

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/services"
	"github.com/mpcarolin/cinematch-server/internal/ui/components"
)

func SearchMovies(c echo.Context) error {
	search := c.QueryParam("search")
	searchQuery := cleanUpSearchQuery(search)
	ctx := c.(*models.RequestContext)

	slog.Info("SearchMovies", "searchQuery", searchQuery)

	response, err := services.SearchMoviesCached(ctx.Cache, searchQuery)
	if err != nil {
		slog.Error("Error fetching movie search results", "error", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	searchResults := []models.Movie{}
	for _, movie := range response.Results {
		if strings.Contains(strings.ToLower(movie.Title), searchQuery) {
			searchResults = append(searchResults, movie)
		}
	}

	component := components.MovieResults(searchResults)

	return component.Render(context.Background(), c.Response().Writer)
}

func cleanUpSearchQuery(search string) string {
	return strings.TrimSpace(strings.ToLower(search))
}