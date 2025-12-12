package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/chat-ai/internal/server/handler"
)

func registerAuthRoutes(r *echo.Group, h *handler.Handlers) {
	authRoute := r.Group("/auth")
	{
		authRoute.POST("/login", h.Auth.LoginHandler())
		authRoute.POST("/register", h.Auth.RegisterHandler())
	}
}
