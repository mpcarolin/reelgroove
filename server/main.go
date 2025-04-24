package main

import (
	"context"
	"log/slog"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/mpcarolin/cinematch-server/internal/components"
	"github.com/mpcarolin/cinematch-server/internal/constants/env"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

func main() {
	e := echo.New()

	// Set debug mode based on environment
	e.Debug = utils.GetEnv() != env.Production

	// Middleware
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.RateLimiter(echoMiddleware.NewRateLimiterMemoryStore(20)))
	e.Use(echoMiddleware.CORSWithConfig(utils.GetCORSConfig()))

	// Routes
	e.GET("/status", func(c echo.Context) error {
		response := "OK " + time.Now().Format("2006-01-02 15:04:05")
		return c.String(http.StatusOK, response)
	})

	e.GET("/home", func(c echo.Context) error {
		component := components.Page(components.MovieSearch())
		return component.Render(context.Background(), c.Response().Writer)
	})


	e.GET("/movie/:movieId/recommendations/:recommendationId", func(c echo.Context) error {
		movieId, err := strconv.Atoi(c.Param("movieId"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid movie id")
		}

		recommendationId, err := strconv.Atoi(c.Param("recommendationId"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid recommendation id")
		}


		recommendations, err := utils.GetMovieRecommendations(movieId)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		nextTrailer, err := utils.GetBestMovieTrailer(recommendationId)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		userLikes, err := utils.GetUserLikesFromCookie(c)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		slog.Info("userLikes", "userLikes", userLikes)

		return components.Page(
			components.Recommendations(models.AppContext{
				MovieId: movieId,
				Trailer: nextTrailer,
				Recommendations: recommendations.Results,
				UserLikes: userLikes,
			}),
		).Render(context.Background(), c.Response().Writer)
	})

	e.PUT("/movie/:movieId/recommendations/:recommendationId/:action", func(c echo.Context) error {
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

		recommendations, err := utils.GetMovieRecommendations(movieId)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		currentRecommendationIndex := slices.IndexFunc(recommendations.Results, func(recommendation utils.Movie) bool { return recommendation.Id == recommendationId })
		nextRecommendationIndex := math.Min(float64(currentRecommendationIndex + 1), float64(len(recommendations.Results) - 1));
		nextRecommendationId := recommendations.Results[int(nextRecommendationIndex)].Id

		nextRecommendationUrl := "/movie/" + strconv.Itoa(movieId) + "/recommendations/" + strconv.Itoa(nextRecommendationId)
		switch action {
		case "skip":
			userLikes = slices.DeleteFunc(userLikes, func(like string) bool { return like == strconv.Itoa(recommendationId) })
			userLikesCookie := utils.CreateUserLikesCookie(userLikes)
			http.SetCookie(c.Response().Writer, userLikesCookie)
			c.Response().Header().Set("HX-Redirect", nextRecommendationUrl)
			return c.NoContent(http.StatusOK);
		case "maybe":
			userLikes = append(userLikes, strconv.Itoa(recommendationId));
			userLikesCookie := utils.CreateUserLikesCookie(userLikes)
			http.SetCookie(c.Response().Writer, userLikesCookie)
			c.Response().Header().Set("HX-Redirect", nextRecommendationUrl)
			return c.NoContent(http.StatusOK);
		case "watch":
			// TODO: implement...
			return nil;
		}
		return nil
	})

	e.GET("/movie/:id/recommendations", func(c echo.Context) error {
		movieId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid movie id")
		}

		recommendations, err := utils.GetMovieRecommendations(movieId)
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
	})

	e.GET("/movies", func(c echo.Context) error {
		search := c.QueryParam("search")
		response, err := utils.SearchMovies(search)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		movies := response.Results
		searchResults := []utils.Movie{}
		for _, movie := range movies {
			if strings.Contains(strings.ToLower(movie.Title), strings.ToLower(search)) {
				searchResults = append(searchResults, movie)
			}
		}

		component := components.MovieResults(searchResults)

		time.Sleep(1 * time.Second) // 1 second delay for testing

		return component.Render(context.Background(), c.Response().Writer)
	})

	e.Static("/assets", "assets")

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
