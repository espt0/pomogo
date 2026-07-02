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

// Authentication
func (h Handler) Registrar(c *echo.Context) error {
	return c.String(http.StatusFound, "Todas as TASKS")
}
func (h Handler) Login(c *echo.Context) error {
	return c.String(http.StatusFound, "TASK específica")
}
func (h Handler) RefreshToken(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h Handler) Logout(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

// Sessions
func (h Handler) IniciaSession(c *echo.Context) error {
	return c.String(http.StatusFound, "TASK específica")
}
func (h Handler) FinalizaSession(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h Handler) Historico(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

// Settings
func (h Handler) VerSettings(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h Handler) AtualizaSettings(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

// Timer
func (h Handler) AtualizaOTimer(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

/*
LEMBRAR
1-Bind do JSON/params/query para struct (DTO de request) — c.Bind(&req)
2-Validação sintática dos dados de entrada — validator, campos obrigatórios, formato, tamanho (validação de "forma", não de regra de negócio)
3-Extrair dados do contexto/autenticação — user_id, claims, tenant_id vindos do middleware (c.Get(...))
4-Chamar o service — repassando DTO + dados de contexto, passando c.Request().Context() adiante
5-Tratar o erro retornado pelo service — mapear erro de domínio (ex: ErrNotFound, ErrConflict, ErrForbidden) para status HTTP correto
6-Montar a struct de resposta (DTO de response) — a partir do que o service retornou, sem vazar model de domínio/banco direto
7-Definir o status code apropriado — 200/201/204/400/401/403/404/409/422/500
8-Retornar a resposta — c.JSON(status, resp)
9-Logging de erro inesperado (opcional) — se algo não mapeado acontecer, logar antes de devolver 500 genérico
*/
