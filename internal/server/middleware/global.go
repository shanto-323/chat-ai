package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shanto-323/chat-ai/internal/server"
)

type Global struct {
	server *server.Server
}

func NewGlobal(s *server.Server) *Global {
	return &Global{
		server: s,
	}
}

func (g *Global) CROS() echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: g.server.Config.Server.CORSAllowedOrigins,
	})
}

