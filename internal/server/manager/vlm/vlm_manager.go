package vlm

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/shanto-323/chat-ai/config"
	"github.com/shanto-323/chat-ai/internal/server/manager/vlm/mock"
	"github.com/shanto-323/chat-ai/model"
)

type VLMManager interface {
	AnalyzeImage(imageURL []string) (*[]model.VLMResult, error)
}

func New(cfg *config.Config,logger *zerolog.Logger) (VLMManager, error) {
	switch cfg.Ai.VLMInterfaceProvider {
	case "mock":
		return mock.NewMockVLM(logger),nil 
	default:
		return nil, fmt.Errorf("so such provider for vlm")
	}
}
