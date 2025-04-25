package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/components"
)

func GetHome(c echo.Context) error {
	component := components.Page(components.MovieSearch())
	return component.Render(context.Background(), c.Response().Writer)
}