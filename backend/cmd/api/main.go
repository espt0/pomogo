package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
)

func main() {
	e := echo.New()

	e.GET("/", teste)

	if err := e.Start(":8181"); err != nil {
		slog.Error("erro ao iniciar servidor", "error", err)
		os.Exit(1)
	}
}

func teste(c *echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
