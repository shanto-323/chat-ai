package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shanto-323/chat-ai/internal/database"
	"github.com/shanto-323/chat-ai/internal/server/errs"
	"github.com/shanto-323/chat-ai/internal/server/manager"
	"github.com/shanto-323/chat-ai/internal/service/image"
	"github.com/shanto-323/chat-ai/model"
	"github.com/shanto-323/chat-ai/model/dto"
	"github.com/shanto-323/chat-ai/model/entity"
)

type ChatService struct {
	is *image.ImageService

	manager *manager.AIManager
	db      database.Database
}

func NewChatService(manager *manager.AIManager, db database.Database, is *image.ImageService) *ChatService {
	return &ChatService{
		manager: manager,
		db:      db,
		is:      is,
	}
}

func (s *ChatService) MultimodalChat(c echo.Context, payload *dto.ChatRequest) (*entity.ConversationLog, error) {
	///-> vlm returns an array of text ...
	response, err := s.manager.LLMManager.GenerateResponse(context.Background(), &dto.LLMRequest{Messages: payload.UserMessage, Model: payload.ModelConfig})
	if err != nil {
		return nil, err
	}

	userId, ok := c.Get("id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	conversationLog := &entity.ConversationLog{
		UserID:       userId,
		TextQuery:    payload.UserMessage,
		ImageURL:     make([]string, 0),
		ResponseText: response,
	}

	conversationLog.LLMModelName = payload.ModelConfig.LLMModel
	conversationLog.VLMModelName = "" /// setup vlm model name

	return s.db.CreateConversationLog(context.Background(), conversationLog)
}

func (s *ChatService) MultimodalChatHistory(c echo.Context, payload *dto.ConversationHistoryQuery) (*model.PaginatedResponse[entity.ConversationLog], error) {
	userId, ok := c.Get("id").(uuid.UUID)
	if !ok {
		return nil, errs.NewInternalServerError()
	}

	return s.db.GetConversationLogHistory(context.Background(), userId, payload)
}
