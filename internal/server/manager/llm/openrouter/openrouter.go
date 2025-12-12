package openrouter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/shanto-323/chat-ai/config"
	"github.com/shanto-323/chat-ai/model/dto"
)

const (
	OPENROUTER_URL = "https://openrouter.ai/api/v1/chat/completions"
	LLAMA_70B      = "meta-llama/llama-3.3-70b-instruct:free"
)

type Openrouter struct {
	Config *config.Config
}

func NewOpenrouter(cfg *config.Config) *Openrouter {
	return &Openrouter{
		Config: cfg,
	}
}

func (o *Openrouter) GenerateResponse(ctx context.Context, request *dto.LLMRequest) (string, error) {
	switch request.Model.LLMModel {
	case "llama_70b":
		return o.response(request.Messages, LLAMA_70B)
	default:
		return "", fmt.Errorf("model not found")
	}
}

func (o *Openrouter) response(query string, model string) (string, error) {
	body := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": query,
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", OPENROUTER_URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	// will add a complition timer
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.Config.Ai.LLMInterfaceApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return o.extractResponse(respBody)
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (o *Openrouter) extractResponse(byteString []byte) (string, error) {
	var resp ChatResponse
	if err := json.Unmarshal(byteString, &resp); err != nil {
		return "", nil
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return resp.Choices[0].Message.Content, nil
}
