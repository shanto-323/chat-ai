package service

import (
	"github.com/shanto-323/chat-ai/internal/server"
	"github.com/shanto-323/chat-ai/internal/service/image"
)

type Services struct {
	Auth AuthService
	Chat ChatService
}

func New(s *server.Server) *Services {
	imageService := image.New(s.Logger)
	return &Services{
		Auth: NewAuthService(s.Config, s.DB),
		Chat: NewChatService(s.Manager, s.DB, imageService),
	}
}
