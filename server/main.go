package main

import (
	"context"
	"net/http"
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

	e.GET("/app", func(c echo.Context) error {
		component := components.Page(components.MovieSearch())
		return component.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/recommendations/:id", func(c echo.Context) error {
		id := c.Param("id")
		component := components.Page(components.Recommendations(id))
		return component.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/movies", func(c echo.Context) error {
		search := c.QueryParam("search")
		movies := []components.Movie{
			{Id: "f47ac10b-58cc-4372-a567-0e02b2c3d479", Title: "Alien", Year: "1979", Poster: "https://m.media-amazon.com/images/M/MV5BY2RlZDYzZjktYjIiMzYwZjYtNDYyZi04YjYzLTg1YmQzZjU4YzI1ZjU@._V1_SX300.jpg"},
			{Id: "7f8d35a2-e34b-4a8c-9f2d-3e6b4c5d8a9f", Title: "The Matrix", Year: "1999", Poster: "https://m.media-amazon.com/images/M/MV5BNzM3NDFhYjctMTNlZS00MTI5LTg0ZTEtNjJkMTZmMWIwYjRmXkEyXkFqcGc@._V1_SX300.jpg"},
			{Id: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d", Title: "The Dark Knight", Year: "2008", Poster: "https://m.media-amazon.com/images/M/MV5BMTMxNTMwODM0NF5BMl5BanBnXkFtZTcwODAyMTk2MzE@._V1_SX300.jpg"},
			{Id: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e", Title: "The Dark Knight Rises", Year: "2012", Poster: "https://m.media-amazon.com/images/M/MV5BMTk4ODQzNDY3Mkg@._V1_SX300.jpg"},
		}
		searchResults := []components.Movie{}
		for _, movie := range movies {
			if strings.Contains(strings.ToLower(movie.Title), strings.ToLower(search)) {
				searchResults = append(searchResults, movie)
			}
		}
		e.Logger.Infof("searchResult Length: %d, search: %s", len(searchResults), search)

		component := components.MovieResults(searchResults)

		time.Sleep(1 * time.Second) // 1 second delay for testing

		return component.Render(context.Background(), c.Response().Writer)
	})

	e.Static("/assets", "assets")

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
