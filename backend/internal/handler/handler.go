package handler

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/espt0/pomogo/internal/apperrors"
	"github.com/espt0/pomogo/internal/auth"
	"github.com/espt0/pomogo/internal/model"
	"github.com/espt0/pomogo/internal/service"

	"github.com/gofrs/uuid/v5"
	"github.com/golang-jwt/jwt/v5"
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

	if err := h.service.CreateUser(ctx, req); err != nil {
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
func (h *Handler) GetCurrentUser(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, err := h.extractUserID(c)
	if err != nil {
		return err
	}

	user, err := h.service.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"erro": "usuário não encontrado"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"erro": "erro interno"})
	}

	return c.JSON(http.StatusOK, user)
}
func (h *Handler) UpdateCurrentUser(c *echo.Context) error {
	ctx := c.Request().Context()
	req := new(model.UpdateUserRequest)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"erro": "corpo da requisição inválido"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"erro": err.Error()})
	}

	userID, err := h.extractUserID(c)
	if err != nil {
		return err
	}

	if err = h.service.UpdateUser(ctx, userID, req); err != nil {
		if errors.Is(err, apperrors.ErrEmailAlreadyExists) {
			return c.JSON(http.StatusConflict, map[string]string{"erro": "email já cadastrado"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"erro": "erro interno"})
	}

	return c.JSON(http.StatusOK, map[string]string{"mensagem": "usuário atualizado com sucesso"})
}
func (h *Handler) UpdatePassword(c *echo.Context) error {
	ctx := c.Request().Context()
	req := new(model.UpdatePasswordUserRequest)

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"erro": "corpo da requisição inválido"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"erro": err.Error()})
	}

	userID, err := h.extractUserID(c)
	if err != nil {
		return err
	}

	if err := h.service.UpdatePasswordUser(ctx, userID, req); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"erro": "erro interno"})
	}

	return c.JSON(http.StatusOK, map[string]string{"mensagem": "password atualizado com sucesso"})
}
func (h *Handler) DeleteCurrentUser(c *echo.Context) error {
	ctx := c.Request().Context()

	userID, err := h.extractUserID(c)
	if err != nil {
		return err
	}

	if err := h.service.DeleteUser(ctx, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"erro": "erro interno"})
	}

	return c.JSON(http.StatusOK, map[string]string{"mensagem": "usuário deletado com sucesso"})
}

// Task
func (h *Handler) CreateTask(c *echo.Context) error {
	return c.JSON(http.StatusCreated, map[string]string{"Task": "CRIADO"})
}
func (h *Handler) ListTasks(c *echo.Context) error {
	return c.String(http.StatusFound, "Todas as TASKS")
}
func (h *Handler) GetTask(c *echo.Context) error {
	return c.String(http.StatusFound, "TASK específica")
}
func (h *Handler) UpdateTask(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h *Handler) DeleteTask(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

// Sessions
func (h *Handler) StartSession(c *echo.Context) error {
	return c.String(http.StatusFound, "TASK específica")
}
func (h *Handler) EndSession(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h *Handler) ListSessions(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

// Settings
func (h *Handler) GetSettings(c *echo.Context) error {
	return c.String(http.StatusFound, "Atualiza TASK")
}
func (h *Handler) UpdateSettings(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

// Timer
func (h *Handler) StreamTimer(c *echo.Context) error {
	return c.String(http.StatusFound, "Deleta TASK")
}

func (h *Handler) extractUserID(c *echo.Context) (uuid.UUID, error) {
	// Pega o token que o middleware colocou no contexto
	token, err := echo.ContextGet[*jwt.Token](c, "user")
	if err != nil {
		return uuid.Nil, echo.ErrUnauthorized
	}

	// Converte os chaims para struct customizada
	claims, ok := token.Claims.(*auth.Claims)
	if !ok {
		return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError, "erro ao processar token")
	}

	return claims.UserID, nil
}
