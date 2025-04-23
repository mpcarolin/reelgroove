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
	});

	e.GET("/app", func(c echo.Context) error {
		component := components.Page()
		return component.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/movies", func(c echo.Context) error {
		search := c.QueryParam("search")
		movies := []components.Movie{
			{Title: "Alien", Year: "1979", Poster: "https://m.media-amazon.com/images/M/MV5BY2RlZDYzZjktYjIiMzYwZjYtNDYyZi04YjYzLTg1YmQzZjU4YzI1ZjU@._V1_SX300.jpg"},
			{Title: "The Matrix", Year: "1999", Poster: "https://m.media-amazon.com/images/M/MV5BNzM3NDFhYjctMTNlZS00MTI5LTg0ZTEtNjJkMTZmMWIwYjRmXkEyXkFqcGc@._V1_SX300.jpg"},
			{Title: "The Dark Knight", Year: "2008", Poster: "https://m.media-amazon.com/images/M/MV5BMTMxNTMwODM0NF5BMl5BanBnXkFtZTcwODAyMTk2MzE@._V1_SX300.jpg"},
			{Title: "The Dark Knight Rises", Year: "2012", Poster: "https://m.media-amazon.com/images/M/MV5BMTk4ODQzNDY3Mkg@._V1_SX300.jpg"},
		}
		searchResults := []components.Movie{}
		for _, movie := range movies {
			if strings.Contains(movie.Title, search) {
				searchResults = append(searchResults, movie)
			}
		}
		component := components.MovieResults(searchResults)
		return component.Render(context.Background(), c.Response().Writer)
	})

	e.Static("/assets", "assets")

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}