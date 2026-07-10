package mw

import (
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func TimeoutMiddleware() middleware.ContextTimeoutConfig {
	return middleware.ContextTimeoutConfig{
		Timeout: 10 * time.Second,
		Skipper: func(c *echo.Context) bool {
			return c.Path() == "/api/v1/time/stream"
		},
	}
}
