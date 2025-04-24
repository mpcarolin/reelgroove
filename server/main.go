package main

import (
	"context"
	"net/http"
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

	// e.POST("/movie/:id/recommendations/:action", func(c echo.Context) error {
	// 	id := c.Param("id")
	// 	action := c.Param("action")

	// 	cookie, err := c.Cookie("recommendation_ids")
	// 	if err != nil {
	// 		return c.String(http.StatusInternalServerError, err.Error())
	// 	}

	// 	recommendation_ids := strings.Split(cookie.Value, ",")

	// 	next_recommendation_id := recommendation_ids[0]

	// 	switch action {
	// 	case "skip":
	// 		return c.String(http.StatusOK, "Skipping recommendation")
	// 	case "maybe":
	// 		return c.String(http.StatusOK, "Maybe recommendation")
	// 	case "watch":
	// 		return c.String(http.StatusOK, "Watching recommendation")
	// 	}
	// 	return nil
	// })

	e.GET("/movie/:id/recommendations", func(c echo.Context) error {
		id := c.Param("id")

		recommendations, err := utils.GetMovieRecommendations(id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		recommendationIds := []string{}
		for _, movie := range recommendations.Results {
			recommendationIds = append(recommendationIds, strconv.Itoa(movie.Id))
		}

		// Create an http cookie to store the recommendation result ids
		cookie := &http.Cookie{
			Name:  "recommendation_ids",
			Value: strings.Join(recommendationIds, ","),
			HttpOnly: true,
			Path:  "/",
		}
		http.SetCookie(c.Response().Writer, cookie)

		trailer, err := utils.GetBestMovieTrailer(id)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		component := components.Page(components.Recommendations(trailer))
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
