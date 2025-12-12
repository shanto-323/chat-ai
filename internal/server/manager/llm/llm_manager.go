package llm

import (
	"context"
	"fmt"

	"github.com/shanto-323/chat-ai/config"
	"github.com/shanto-323/chat-ai/internal/server/manager/llm/openrouter"
	"github.com/shanto-323/chat-ai/model/dto"
)

type LLMManager interface {
	GenerateResponse(ctx context.Context, request *dto.LLMRequest) (string, error)
}

func New(cfg *config.Config) (LLMManager, error) {
	switch cfg.Ai.LLMInterfaceProvider {
	case "openrouter":
		return openrouter.NewOpenrouter(cfg), nil
	default:
		return nil, fmt.Errorf("so such provider")
	}
}
