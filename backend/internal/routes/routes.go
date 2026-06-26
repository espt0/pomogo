package routes

import (
	"github.com/espt0/pomogo/internal/handler"

	"github.com/labstack/echo/v5"
)

// Criar, Ler, Atualizar, Excluir (CRUD)
func SetupRoutes(e *echo.Echo, h *handler.Handler) error {

	v1 := e.Group("/api/v1")

	protected := v1.Group("")

	// Users
	protected.GET("/users/me", h.ListaUser)
	protected.PATCH("/users/me", h.AtualizaUser)
	protected.PATCH("/users/me/password", h.AtualizaSenha)
	protected.DELETE("/users/me", h.DeletaUser)

	// Task
	protected.POST("/task", h.CriaTask)
	protected.GET("/task", h.ListaTodasTasks)
	protected.GET("/task/:id", h.ListaTaskID)
	protected.PATCH("/task/:id", h.AtualizaTask)
	protected.DELETE("/task/:id", h.DeletaTask)

	return nil
}
