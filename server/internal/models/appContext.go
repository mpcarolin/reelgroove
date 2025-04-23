package models

import (
	"github.com/labstack/echo/v4"
)

type AppContext struct {
	echo.Context
}

type AppHandler func(c *AppContext) error

// ToHandler converts an AppHandler to an echo.HandlerFunc
func (h AppHandler) ToHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		appCtx := c.(*AppContext)
		return h(appCtx)
	}
}