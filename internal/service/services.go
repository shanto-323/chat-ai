package service

import (
	"github.com/shanto-323/chat-ai/internal/server"
)

type Services struct {
}

func New(s *server.Server) *Services {
	return &Services{}
}
