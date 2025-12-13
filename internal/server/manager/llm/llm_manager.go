package llm

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/shanto-323/chat-ai/config"
	"github.com/shanto-323/chat-ai/internal/server/manager/llm/openrouter"
	"github.com/shanto-323/chat-ai/model/dto"
)

type LLMManager interface {
	GenerateResponse(ctx context.Context, request *dto.LLMRequest) (string, error)
}

func New(cfg *config.Config, l *zerolog.Logger) (LLMManager, error) {
	switch cfg.Ai.LLMInterfaceProvider {
	case "openrouter":
		return openrouter.NewOpenrouter(cfg,l), nil
	default:
		return nil, fmt.Errorf("so such provider")
	}
}
