package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/mpcarolin/cinematch-server/internal/models"
	"github.com/mpcarolin/cinematch-server/internal/utils"
)

// RequestContext echo middleware
func SetupRequestContext(next echo.HandlerFunc) echo.HandlerFunc {
	cache := utils.GetCache()
	return func(c echo.Context) error {
		ctx := &models.RequestContext{
			Context: c,
			Cache: cache,
		}
		return next(ctx)
	}
}
