package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/services"
	"github.com/mpcarolin/cinematch-server/internal/ui/components"
)

// Renders the recommendations page for a single movie in the recommendations list
// GET /movie/:movieId/recommendations/:recommendationId/watchproviders
func GetWatchProviders(c echo.Context) error {
	ctx := c.(*models.RequestContext)

	recommendationId, err := strconv.Atoi(c.Param("recommendationId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid recommendation id")
	}

	watchProviders, err := services.GetWatchProvidersCached(ctx.Cache, recommendationId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error fetching watch providers")
	}

	renderCtx := context.Background()	
	components.WatchProvidersAllOptions(watchProviders).Render(renderCtx, c.Response().Writer)

	return nil
}

