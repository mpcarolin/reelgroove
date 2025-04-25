package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/components"
)

func GetAbout(c echo.Context) error {
	component := components.Page(components.About())
	return component.Render(context.Background(), c.Response().Writer)
}
