package handlers

import (
	"context"
	"math"
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/services"
	"github.com/mpcarolin/cinematch-server/internal/ui/pages"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

// Updates the user's likes cookie to add the current recommendation, and re-renders the recommendation page
func LikeRecommendation(c echo.Context) error {
	recommendationViewModel, err := InitRecommendationViewModel(c)
	if err != nil {
		return err
	}

	if recommendationViewModel.NextRecommendationId == recommendationViewModel.CurrentRecommendationId {
		// there are no more recommendations, so redirect to the summary page
		return redirectToSummary(c, recommendationViewModel.MovieId)
	}

	// add the current recommendation to the user's likes cookie, then update the cookie
	recommendationViewModel.UserLikes = append(recommendationViewModel.UserLikes, strconv.Itoa(recommendationViewModel.CurrentRecommendationId))
	userLikesCookie := utils.CreateUserLikesCookie(recommendationViewModel.UserLikes)
	http.SetCookie(c.Response().Writer, userLikesCookie)


	// update the current recommendation id to the next recommendation id
	recommendationViewModel.CurrentRecommendationId = recommendationViewModel.NextRecommendationId
	recommendationViewModel.NextRecommendationId = GetNextRecommendationId(recommendationViewModel.Recommendations, recommendationViewModel.NextRecommendationId)

	// update the client url to the next recommendation
	nextRecommendationUrl := models.GetRecommendationUrl(recommendationViewModel.MovieId, recommendationViewModel.NextRecommendationId, &recommendationViewModel.Settings.Autoplay)

	c.Response().Header().Set(
		"HX-Push-Url",
		nextRecommendationUrl,
	)

	return pages.Recommendation(recommendationViewModel).Render(context.Background(), c.Response().Writer)
}

// Updates the user's likes cookie to remove the current recommendation, if it was in the cookie, and re-renders the recommendation page
func SkipRecommendation(c echo.Context) error {
	recommendationViewModel, err := InitRecommendationViewModel(c)
	if err != nil {
		return err
	}

	if recommendationViewModel.NextRecommendationId == recommendationViewModel.CurrentRecommendationId {
		// there are no more recommendations, so redirect to the summary page
		return redirectToSummary(c, recommendationViewModel.MovieId)
	}

	recommendationViewModel.UserLikes = slices.DeleteFunc(recommendationViewModel.UserLikes, func(like string) bool { return like == strconv.Itoa(recommendationViewModel.CurrentRecommendationId) })
	userLikesCookie := utils.CreateUserLikesCookie(recommendationViewModel.UserLikes)
	http.SetCookie(c.Response().Writer, userLikesCookie)

	// update the current recommendation id to the next recommendation id
	recommendationViewModel.CurrentRecommendationId = recommendationViewModel.NextRecommendationId
	recommendationViewModel.NextRecommendationId = GetNextRecommendationId(recommendationViewModel.Recommendations, recommendationViewModel.NextRecommendationId)

	// update the client url to the next recommendation
	c.Response().Header().Set(
		"HX-Push-Url",
		models.GetRecommendationUrl(recommendationViewModel.MovieId, recommendationViewModel.NextRecommendationId, &recommendationViewModel.Settings.Autoplay),
	)

	return pages.Recommendation(recommendationViewModel).Render(context.Background(), c.Response().Writer)
}

func redirectToSummary(c echo.Context, movieId int) error {
	summaryUrl := "/movie/" + strconv.Itoa(movieId) + "/recommendation-summary"

	c.Response().Header().Set(
		"HX-Redirect",
		summaryUrl,
	)

	return c.String(http.StatusOK, "Redirecting to summary page")
}

// TODO: refactor. This is slightly confusing because it gets the NEXT recommendation id's trailer!
// This idea of current and next recommendation ids is a little confusing.
func InitRecommendationViewModel(c echo.Context) (*pages.RecommendationViewModel, error) {
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

	nextRecommendationId := GetNextRecommendationId(recommendations.Results, recommendationId)

	nextTrailer, err := services.GetBestMovieTrailerCached(ctx.Cache, nextRecommendationId)
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, err.Error())
	}

	recommendationViewModel := pages.RecommendationViewModel{
		MovieId:                 movieId,
		CurrentRecommendationId: recommendationId,
		NextRecommendationId:    nextRecommendationId,
		Recommendations:         recommendations.Results,
		UserLikes:               userLikes,
		Trailer:                 nextTrailer,
		Settings:                NewRecommendationSettings(c),
	}

	return &recommendationViewModel, nil
}

func IsLastRecommendation(recommendations []models.Movie, currentRecommendationId int) bool {
	currentRecommendationIndex := slices.IndexFunc(recommendations, func(recommendation models.Movie) bool { return recommendation.Id == currentRecommendationId })
	return currentRecommendationIndex == len(recommendations) - 1
}

// Returns the next recommendation id, or -1 if the current recommendation is the last one
func GetNextRecommendationId(recommendations []models.Movie, currentRecommendationId int) int {
	currentRecommendationIndex := slices.IndexFunc(recommendations, func(recommendation models.Movie) bool { return recommendation.Id == currentRecommendationId })
	nextRecommendationIndex := math.Min(float64(currentRecommendationIndex+1), float64(len(recommendations)-1))
	return recommendations[int(nextRecommendationIndex)].Id
}

func NewRecommendationSettings(c echo.Context) models.RecommendationSettings {
	query := c.QueryParam("autoplay")
	formValue := c.FormValue("autoplay")
	clientUrl, _ := GetClientUrl(c)

	// try to fetch the autoplay setting in order of priority:
	// 1. query param
	// 2. form data
	// 3. fall back to existing client url value
	autoplay := false
	if query != "" {
		autoplay = query == "on"
	} else if formValue != "" {
		autoplay = formValue == "on"
	} else if clientUrl != nil {
		autoplay = clientUrl.Query().Get("autoplay") == "on"
	}

	settings := models.RecommendationSettings{
		Autoplay: autoplay,
	}

	return settings
}

func NewRecommendationSettingsFromForm(c echo.Context) models.RecommendationSettings {
	formValue := c.FormValue("autoplay")
	autoplay := formValue == "on"

	return models.RecommendationSettings{
		Autoplay: autoplay,
	}
}

func NewRecommendationSettingsFromQuery(c echo.Context) models.RecommendationSettings {
	query := c.QueryParam("autoplay")
	autoplay := query == "on"

	return models.RecommendationSettings{
		Autoplay: autoplay,
	}
}

func NewRecommendationSettingsFromClientUrl(c echo.Context) models.RecommendationSettings {
	clientUrl, _ := GetClientUrl(c)
	autoplay := clientUrl.Query().Get("autoplay") == "on"

	return models.RecommendationSettings{
		Autoplay: autoplay,
	}
}

func GetClientUrl(c echo.Context) (*url.URL, error) {
	currentUrl := c.Request().Header.Get("Hx-Current-Url")
	return url.Parse(currentUrl)
}
