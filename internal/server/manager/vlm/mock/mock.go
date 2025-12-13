package mock

import (
	"github.com/rs/zerolog"
	"github.com/shanto-323/chat-ai/model"
)

type MockVLM struct {
	logger *zerolog.Logger
}

func NewMockVLM(l *zerolog.Logger) *MockVLM {
	return &MockVLM{
		logger: l,
	}
}

func (mvs *MockVLM) AnalyzeImage(imageURL []string) (*[]model.VLMResult, error) {
	mvs.logger.Info().Msg("Image analyzing done!!")
	return nil, nil
}
