package handler

import (
	"net/http"

	"github.com/espt0/pomogo/internal/service"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// Task
func (h Handler) CriaTask(c *echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]string{"Task": "CRIADO"})
}
func (h Handler) ListaTodasTasks(c *echo.Context) error {
	return c.String(http.StatusFound, "Todas as TASKS")
}
func (h Handler) ListaTaskID(c *echo.Context) error {
	return c.String(http.StatusFound, "TASK específica")
}
func (h Handler) AtualizaTask(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h Handler) DeletaTask(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

// User
func (h Handler) ListaUser(c *echo.Context) error {
	return c.String(http.StatusFound, "Todas as TASKS")
}
func (h Handler) AtualizaUser(c *echo.Context) error {
	return c.String(http.StatusFound, "TASK específica")
}
func (h Handler) AtualizaSenha(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h Handler) DeletaUser(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}
