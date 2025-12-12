package manager

import (
	"github.com/shanto-323/chat-ai/config"
	"github.com/shanto-323/chat-ai/internal/server/manager/llm"
	"github.com/shanto-323/chat-ai/internal/server/manager/vlm"
)

type AIManager struct {
	LLMManager llm.LLMManager
	VLMManager vlm.VLMManager
}

func New(cfg *config.Config) (*AIManager, error) {
	llmManager, err := llm.New(cfg)
	if err != nil {
		return nil, err
	}

	vlmManager, err := vlm.New(cfg)
	if err != nil {
		return nil, err
	}

	return &AIManager{
		LLMManager: llmManager,
		VLMManager: vlmManager,
	}, nil
}
