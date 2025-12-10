package dto

import (
	"github.com/go-playground/validator"
	"github.com/shanto-323/chat-ai/model"
)

type ConversationLogRequest struct {
	TextQuery string   `json:"text_query,omitempty"`
	ImageURL  []string `json:"image_url,omitempty"`
}

func (l *ConversationLogRequest) Validate() error {
	return validator.New().Struct(l)
}

type ConversationHistoryQuery struct {
	Page  *int    `query:"page" validate:"omitempty,min=1"`
	Limit *int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Order *string `query:"order" validate:"omitempty,oneof=asc desc"`
}

func (l *ConversationHistoryQuery) Validate() error {
	if err := validator.New().Struct(l); err != nil {
		return err
	}

	if l.Page == nil {
		defaultPage := 1
		l.Page = &defaultPage
	}
	if l.Limit == nil {
		defaultLimit := 20
		l.Limit = &defaultLimit
	}
	if l.Order == nil {
		defaultOrder := "desc"
		l.Order = &defaultOrder
	}

	return nil
}

type ConversationLogResponse struct {
	model.BaseLV
	TextQuery    string   `json:"text_query"`
	ImageURL     []string `json:"image_url"`
	ResponseText string   `json:"response_text"`
}
