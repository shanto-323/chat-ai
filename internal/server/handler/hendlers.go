package handler

import (
	"github.com/shanto-323/chat-ai/internal/server"
	"github.com/shanto-323/chat-ai/internal/service"
)

type Handlers struct {
}

func New(s *server.Server, services *service.Services) *Handlers {
	return &Handlers{}
}
