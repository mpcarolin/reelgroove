package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/ui"
	"github.com/mpcarolin/cinematch-server/internal/ui/pages"
)

func GetAbout(c echo.Context) error {
	component := ui.Page(pages.About())
	return component.Render(context.Background(), c.Response().Writer)
}
