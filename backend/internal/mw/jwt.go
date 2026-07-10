package mw

import (
	"net/http"
	"os"

	"github.com/espt0/pomogo/internal/auth"
	"github.com/golang-jwt/jwt/v5"

	echojwt "github.com/labstack/echo-jwt/v5"
	"github.com/labstack/echo/v5"
)

func JWTMiddleware() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
		NewClaimsFunc: func(c *echo.Context) jwt.Claims {
			return new(auth.Claims)
		},
		ErrorHandler: func(c *echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, "token inválido ou ausente")
		},
	})
}
