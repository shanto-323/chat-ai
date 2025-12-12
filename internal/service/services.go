package service

import (
	"github.com/shanto-323/chat-ai/internal/server"
)

type Services struct {
	Auth *AuthService
	Chat *ChatService
}

func New(s *server.Server) *Services {
	return &Services{
		Auth: NewAuthService(s.Config, s.DB),
		Chat: NewChatService(s.Manager, s.DB),
	}
}
