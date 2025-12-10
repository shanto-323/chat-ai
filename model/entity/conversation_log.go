package entity

import (
	"github.com/google/uuid"
	"github.com/shanto-323/chat-ai/model"
)

type ConversationLog struct {
	model.BaseId
	model.BaseLV

	UserID       uuid.UUID `db:"user_id"`
	TextQuery    string    `db:"text_query"`
	ImageURL     []string  `db:"image_url"`
	ResponseText string    `db:"response_text"`
}
