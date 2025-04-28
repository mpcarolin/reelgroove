package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/mpcarolin/cinematch-server/internal/constants/env"
	"github.com/mpcarolin/cinematch-server/internal/handlers"
	"github.com/mpcarolin/cinematch-server/internal/middleware"
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

	e.Use(middleware.SetupRequestContext)

	// Routes
	e.GET("/status", func(c echo.Context) error {
		response := "OK " + time.Now().Format("2006-01-02 15:04:05")
		return c.String(http.StatusOK, response)
	})

	e.GET("/about", handlers.GetAbout)
	e.GET("/home", handlers.GetHome)
	e.GET("/movies", handlers.SearchMovies)
	e.GET("/movie/:movieId/recommendations", handlers.GetRecommendations)
	e.GET("/movie/:movieId/recommendations/:recommendationId", handlers.GetRecommendationById)
	e.PUT("/movie/:movieId/recommendations/:recommendationId/:action", handlers.HandleRecommendationAction)
	e.Static("/assets", "assets")

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
