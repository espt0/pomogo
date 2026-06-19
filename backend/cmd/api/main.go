package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/espt0/pomogo/internal/db"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("erro ao carregar o arquivo .env", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	DB, err := db.Connect()
	if err != nil {
		slog.Error("erro ao conectar ao banco de dados", "error", err)
		os.Exit(1)
	}
	defer DB.Close()

	err = DB.Ping(ctx)
	if err != nil {
		slog.Error("erro ao pingar banco de dados", "error", err)
		os.Exit(1)
	}

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

/*
1. Conexão com banco ✅
2. Ping no banco ✅
3. Configurar rotas/handlers ✅
4. defer db.Close() ✅
5. e.Start() (por último) ✅
*/
