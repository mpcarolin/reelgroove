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

func UpdateRecommendationSettings(c echo.Context) error {
	ctx := c.(*models.RequestContext)
	// TODO: should we make functions for these? Across error handlers they are
	// repeating the same messages and codes, but maybe the duplication is worth it over extra abstraction...
	movieId, err := strconv.Atoi(c.Param("movieId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid movie id")
	}

	recommendationId, err := strconv.Atoi(c.Param("recommendationId"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid recommendation id")
	}

	trailer, err := services.GetBestMovieTrailerCached(ctx.Cache, recommendationId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Error getting trailer")
	}

	settings := NewRecommendationSettingsFromForm(c)

	// set header first
	c.Response().Header().Set("HX-Push-Url", models.GetRecommendationUrl(movieId, recommendationId, &settings.Autoplay))

	renderCtx := context.Background()

	err = components.YouTubeVideoEmbed(
		trailer.Key,
		components.VideoConfig{Autoplay: settings.Autoplay, OOB: true},
	).Render(renderCtx, c.Response().Writer)
	if err != nil {
		return err
	}

	err = components.TrailerSettings(components.TrailerSettingsViewModel{
		Settings:          settings,
		UpdateSettingsUrl: models.GetUpdateSettingsUrl(movieId, recommendationId),
	}).Render(renderCtx, c.Response().Writer)
	if err != nil {
		return err
	}

	return nil
}
