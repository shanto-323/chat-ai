package middleware

import "github.com/shanto-323/chat-ai/internal/server"

type Middlewares struct {
	*Global
	*RateLimit
	*ContextEnhancer
}

func New(s *server.Server) *Middlewares {
	return &Middlewares{
		Global:          NewGlobal(s),
		RateLimit:       NewRateLimit(s),
		ContextEnhancer: NewContextEnhancer(s),
	}
}
