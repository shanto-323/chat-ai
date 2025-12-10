package model

import "time"

type BaseLLMModel struct {
	LLMModelName string `json:"llmModelName" db:"llm_model_name"`
}

type BaseVLMModel struct {
	VLMModelName string `json:"vlmModelName" db:"vlm_model_name"`
}

type BaseTimestamp struct {
	Timestamp time.Time `json:"timestamp" db:"timestamp"`
}

type BaseLV struct {
	BaseLLMModel
	BaseVLMModel
	BaseTimestamp
}
