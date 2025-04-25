package models

import (
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/labstack/echo/v4"
)

type RequestContext struct {
	echo.Context
	Cache *cache.Cache[string]
}
