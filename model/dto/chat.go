package dto

import (
	"fmt"

	"github.com/go-playground/validator"
)

type ImageData struct {
	URL    string `json:"url,omitempty"`
	Base64 string `json:"base64,omitempty"`
	Type   string `json:"type,omitempty"`
}

type ModelConfig struct {
	LLMModel string `json:"llm_model"`
}

type ChatRequest struct {
	UserMessage string       `json:"message_query,omitempty"`
	Images      []ImageData  `json:"image_s,omitempty"`
	ModelConfig *ModelConfig `json:"model_config,omitempty"`
}

func (r *ChatRequest) Validate() error {
	if err := validator.New().Struct(r); err != nil {
		return err
	}

	if r.UserMessage == "" && len(r.Images) == 0 {
		return fmt.Errorf("atleast one field required")
	}

	if r.ModelConfig == nil {
		r.ModelConfig = &ModelConfig{
			LLMModel: "llama_70b",
		}
	}

	return nil
}
