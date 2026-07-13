package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/espt0/pomogo/internal/db"
	"github.com/espt0/pomogo/internal/handler"
	"github.com/espt0/pomogo/internal/mw"
	"github.com/espt0/pomogo/internal/repository"
	"github.com/espt0/pomogo/internal/routes"
	"github.com/espt0/pomogo/internal/service"
	"github.com/espt0/pomogo/internal/validator"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	ctx := context.Background()

	// Banco de dados
	if err := godotenv.Load("../../.env"); err != nil {
		slog.Error("erro ao carregar o arquivo .env", "error", err)
		os.Exit(1)
	}
	DB, err := db.Connect(ctx)
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

	// Cria instância do Echo
	e := echo.New()

	// Middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLogger())
	// CORS vou add depois
	e.Use(middleware.ContextTimeoutWithConfig(mw.TimeoutMiddleware()))
	e.Use(middleware.Gzip())

	// Validação
	e.Validator = validator.NewEchoValidator()

	// Injeção de dependência
	repo := repository.NewRepository(DB)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)
	err = routes.SetupRoutes(e, handler)
	if err != nil {
		slog.Error("erro ao iniciar rotas", "error", err)
	}

	// Inicia o servidor
	if err := e.Start(":8181"); err != nil {
		slog.Error("erro ao iniciar servidor", "error", err)
		os.Exit(1)
	}
}
