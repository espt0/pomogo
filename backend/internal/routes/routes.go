package routes

import (
	"github.com/espt0/pomogo/internal/handler"
	"github.com/espt0/pomogo/internal/mw"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func SetupRoutes(e *echo.Echo, h *handler.Handler) error {
	v1 := e.Group("/api/v1")

	// Público
	auth := v1.Group("/auth", middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(3.0)))
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.RefreshToken)
	auth.POST("/logout", h.Logout)

	// Protegido
	protected := v1.Group("", mw.JWTMiddleware())

	// Users
	protected.GET("/users/me", h.GetCurrentUser)
	protected.PATCH("/users/me", h.UpdateCurrentUser)
	protected.PATCH("/users/me/password", h.UpdatePassword)
	protected.DELETE("/users/me", h.DeleteCurrentUser)

	// Task
	protected.POST("/task", h.CriaTask)
	protected.GET("/task", h.ListaTodasTasks)
	protected.GET("/task/:id", h.ListaTaskID)
	protected.PATCH("/task/:id", h.AtualizaTask)
	protected.DELETE("/task/:id", h.DeletaTask)

	// Sessions
	protected.POST("/sessions", h.IniciaSession)
	protected.PATCH("/sessions/:id/end", h.FinalizaSession)
	protected.GET("/sessions", h.Historico)

	// Settings
	protected.GET("/settings", h.VerSettings)
	protected.PUT("/settings", h.AtualizaSettings)

	// Timer (SSE)
	protected.GET("/time/stream", h.AtualizaOTimer)

	return nil
}
