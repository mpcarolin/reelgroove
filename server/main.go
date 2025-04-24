package main

import (
	"context"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/mpcarolin/cinematch-server/internal/components"
	"github.com/mpcarolin/cinematch-server/internal/constants/env"
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

	e.PUT("/movie/:movieId/recommendations/:action", func(c echo.Context) error {
		movieId, err := strconv.Atoi(c.Param("movieId"))
		if err != nil {
			return c.String(http.StatusBadRequest, "Invalid movie id")
		}

		action := c.Param("action")

		recommendationIds, err := utils.GetRecommendationIdsFromCookie(c)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		currentRecommendationId, err := utils.GetCurrentRecommendationMovieIdFromCookie(c)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		userLikes, err := utils.GetUserLikesFromCookie(c)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}


		switch action {
		case "skip":
			idx := slices.Index(recommendationIds, strconv.Itoa(currentRecommendationId))
			if idx == -1 {
				// TODO: handle this better
				return c.String(http.StatusInternalServerError, "Current recommendation not found")
			}
			nextRecommendationId := recommendationIds[idx+1]
			nextRecommendationCookie := utils.CreateCurrentRecommendationCookie(nextRecommendationId)
			http.SetCookie(c.Response().Writer, nextRecommendationCookie)

			nextRecommendationIdInt, err := strconv.Atoi(nextRecommendationId)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			recommendations, err := utils.GetMovieRecommendations(movieId)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			// TODO: major problem. We don't want to fetch all these recommendations every time we run this endpoint.
			// We need to preserve them somehow, but cookies aren't the right tool, because of limited space.
			// THis will work for now because this is mocked out function but once it's not, we need to fix this.
			nextTrailer, err := utils.GetBestMovieTrailer(nextRecommendationIdInt)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			return components.Recommendations(nextTrailer, recommendations.Results).Render(context.Background(), c.Response().Writer)
		case "maybe":
			idx := slices.Index(recommendationIds, strconv.Itoa(currentRecommendationId))
			if idx == -1 {
				// TODO: handle this better
				return c.String(http.StatusInternalServerError, "Current recommendation not found")
			}
			nextRecommendationId := recommendationIds[idx+1]
			nextRecommendationCookie := utils.CreateCurrentRecommendationCookie(nextRecommendationId)
			http.SetCookie(c.Response().Writer, nextRecommendationCookie)

			nextRecommendationIdInt, err := strconv.Atoi(nextRecommendationId)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			recommendations, err := utils.GetMovieRecommendations(movieId)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			// TODO: major problem. We don't want to fetch all these recommendations every time we run this endpoint.
			// We need to preserve them somehow, but cookies aren't the right tool, because of limited space.
			// THis will work for now because this is mocked out function but once it's not, we need to fix this.
			nextTrailer, err := utils.GetBestMovieTrailer(nextRecommendationIdInt)
			if err != nil {
				return c.String(http.StatusInternalServerError, err.Error())
			}

			nextUserLikes := append(userLikes, strconv.Itoa(currentRecommendationId))
			nextUserLikesCookie := utils.CreateUserLikesCookie(nextUserLikes)
			http.SetCookie(c.Response().Writer, nextUserLikesCookie)

			return components.Recommendations(nextTrailer, recommendations.Results).Render(context.Background(), c.Response().Writer)
		case "watch":
			// movie, err := utils.GetMovie(currentRecommendationIdInt)
			// if err != nil {
			// 	return c.String(http.StatusInternalServerError, err.Error())
			// }
			// return components.Page(components.Watch(currentRecommendationId, recommendations.Results)).Render(context.Background(), c.Response().Writer)
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

		recommendationIds := []string{}
		for _, movie := range recommendations.Results {
			recommendationIds = append(recommendationIds, strconv.Itoa(movie.Id))
		}

		// Create cookies for recommendation tracking, if one doesn't already exist
		if existingRecommendationIds, _ := utils.GetRecommendationIdsFromCookie(c); len(existingRecommendationIds) == 0 {
			recommendationIdsCookie := utils.CreateRecommendationIdsCookie(recommendationIds)
			http.SetCookie(c.Response().Writer, recommendationIdsCookie)
		}

		currentRecommendationId := recommendationIds[0]
		currentRecommendationCookie := utils.CreateCurrentRecommendationCookie(currentRecommendationId)
		http.SetCookie(c.Response().Writer, currentRecommendationCookie)

		userLikesCookie := utils.CreateUserLikesCookie([]string{})
		http.SetCookie(c.Response().Writer, userLikesCookie)

		currentRecommendationIdInt, err := strconv.Atoi(currentRecommendationId)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		trailer, err := utils.GetBestMovieTrailer(currentRecommendationIdInt)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		component := components.Page(components.Recommendations(trailer, recommendations.Results))
		return component.Render(context.Background(), c.Response().Writer)
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
