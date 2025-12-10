package vlm

import (
	"github.com/shanto-323/chat-ai/config"
	"github.com/shanto-323/chat-ai/model"
)

type VLMManager interface {
	AnalyzeImage(imageURL string) (*model.VLMResult, error)
}

func NewVLMManager(config *config.Config)

