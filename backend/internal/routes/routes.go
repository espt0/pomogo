package routes

import (
	"github.com/espt0/pomogo/internal/handler"
	"github.com/espt0/pomogo/internal/mw"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func SetupRoutes(e *echo.Echo, h *handler.Handler) {
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
	protected.POST("/tasks", h.CreateTask)
	protected.GET("/tasks", h.ListTasks)
	protected.GET("/tasks/:id", h.GetTask)
	protected.PATCH("/tasks/:id", h.UpdateTask)
	protected.DELETE("/tasks/:id", h.DeleteTask)

	// Sessions
	protected.POST("/sessions", h.StartSession)
	protected.PATCH("/sessions/:id/end", h.EndSession)
	protected.GET("/sessions", h.ListSessions)

	// Settings
	protected.GET("/settings", h.GetSettings)
	protected.PUT("/settings", h.UpdateSettings)

	// Timer (SSE)
	protected.GET("/time/stream", h.StreamTimer)

}
