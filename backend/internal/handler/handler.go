package handler

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/espt0/pomogo/internal/apperrors"
	"github.com/espt0/pomogo/internal/model"
	"github.com/espt0/pomogo/internal/service"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

// Autenticação
func (h *Handler) Register(c *echo.Context) error {
	ctx := c.Request().Context()
	req := new(model.CreateUserRequest)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"erro": "corpo da requisição inválido"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"erro": err.Error()})
	}

	err := h.service.CreateUser(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrEmailAlreadyExists) {
			return c.JSON(http.StatusConflict, map[string]string{"erro": "email já cadastrado"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"erro": "erro ao criar usuário"})
	}

	return c.JSON(http.StatusCreated, map[string]string{"mensagem": "usuário criado com sucesso"})
}
func (h *Handler) Login(c *echo.Context) error {
	ctx := c.Request().Context()
	req := new(model.LoginRequest)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"erro": "corpo da requisição inválido"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"erro": err.Error()})
	}

	accessToken, refreshToken, err := h.service.LoginUser(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidEmailOrPassword) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"erro": "email ou senha inválidos"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"erro": "erro interno"})
	}

	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = refreshToken
	cookie.Path = "/api/v1/auth/refresh"
	cookie.HttpOnly = true
	cookie.Secure = os.Getenv("ENV") == "production"
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Expires = time.Now().Add(43200 * time.Minute)

	c.SetCookie(cookie)

	response := model.TokenResponse{
		AccessToken: accessToken,
		TokenType:   "Bearer",
		ExpiresIn:   900,
	}

	return c.JSON(http.StatusOK, response)
}
func (h *Handler) RefreshToken(c *echo.Context) error {
	ctx := c.Request().Context()

	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"erro": "cookie não encontrado"})
	}

	newAccessToken, newRefreshToken, err := h.service.RotateTokens(ctx, cookie.Value)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidRefreshToken) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"erro": "token inválido ou já utilizado"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"erro": "erro ao renovar sessão"})
	}

	cookie = new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = newRefreshToken
	cookie.Path = "/api/v1/auth/refresh"
	cookie.HttpOnly = true
	cookie.Secure = os.Getenv("ENV") == "production"
	cookie.SameSite = http.SameSiteStrictMode
	cookie.Expires = time.Now().Add(43200 * time.Minute)

	c.SetCookie(cookie)

	response := model.TokenResponse{
		AccessToken: newAccessToken,
		TokenType:   "Bearer",
		ExpiresIn:   900,
	}

	return c.JSON(http.StatusOK, response)
}
func (h *Handler) Logout(c *echo.Context) error {
	ctx := c.Request().Context()

	cookie, err := c.Cookie("refreshToken")
	if err == nil {
		_ = h.service.LogoutUser(ctx, cookie.Value)
	}

	newCookie := new(http.Cookie)
	newCookie.Name = "refreshToken"
	newCookie.Value = ""
	newCookie.Path = "/api/v1/auth/refresh"
	newCookie.MaxAge = -1
	newCookie.HttpOnly = true
	newCookie.Secure = os.Getenv("ENV") == "production"
	newCookie.SameSite = http.SameSiteStrictMode

	c.SetCookie(newCookie)

	return c.JSON(http.StatusOK, map[string]string{"mensagem": "usuário deslogado com sucesso"})
}

// User
func (h *Handler) ListaUser(c *echo.Context) error {
	return c.String(http.StatusFound, "Todas as TASKS")
}
func (h *Handler) AtualizaUser(c *echo.Context) error {
	return c.String(http.StatusFound, "TASK específica")
}
func (h *Handler) AtualizaSenha(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h *Handler) DeletaUser(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

// Task
func (h *Handler) CriaTask(c *echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]string{"Task": "CRIADO"})
}
func (h *Handler) ListaTodasTasks(c *echo.Context) error {
	return c.String(http.StatusFound, "Todas as TASKS")
}
func (h *Handler) ListaTaskID(c *echo.Context) error {
	return c.String(http.StatusFound, "TASK específica")
}
func (h *Handler) AtualizaTask(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h *Handler) DeletaTask(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

// Sessions
func (h *Handler) IniciaSession(c *echo.Context) error {
	return c.String(http.StatusFound, "TASK específica")
}
func (h *Handler) FinalizaSession(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h *Handler) Historico(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

// Settings
func (h *Handler) VerSettings(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h *Handler) AtualizaSettings(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

// Timer
func (h *Handler) AtualizaOTimer(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}
